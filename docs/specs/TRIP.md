# Trip Module Specification

**Version:** 1.1  
**Status:** Canonical — implementation source of truth  
**Base path (production):** `/api/trip`  
**Base path (e2e tests):** `/` (module mounted at root)

Do not implement trip behavior from archived product docs. Follow this file.

---

## 1. Domain model

| Entity | Purpose | Editable after handoff? |
|--------|---------|-------------------------|
| **TripRequest** | Customer booking intent + matching lifecycle | **No** — frozen on driver accept |
| **Trip** | Driver assignment + ride execution | Yes — all post-accept lifecycle |

**Handoff rule:** On driver accept → set `trip_request` to `DRIVER_ACCEPTED` (terminal), create `trip` with `TRIP_ACCEPTED`. After that, **only update `trip`** for start, cancel, and complete.

### Entity fields

**TripRequest** (`trip.trip_requests`)

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `customer_id` | UUID | FK → `trip.customers` |
| `origin` | string | Pickup location |
| `destination` | string | Drop-off location |
| `status` | int | See §2 |
| `created_at`, `updated_at` | timestamptz | |

**Trip** (`trip.trips`)

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `trip_request_id` | UUID | FK → `trip.trip_requests`, **unique** (one trip per request) |
| `customer_id` | UUID | FK → `trip.customers`; copied from `trip_request` at accept |
| `driver_id` | UUID | FK → `trip.drivers` |
| `status` | int | See §3 |
| `created_at`, `updated_at` | timestamptz | |

`customer_id` is stored on `trips` at accept time so post-accept ownership checks use the trip row directly.

### Authorization (requirements)

| Action | Who may perform it |
|--------|-------------------|
| Create / read / cancel trip request (pre-accept) | Owning customer only |
| List open trip requests | Any authenticated driver |
| Accept trip request | Any authenticated driver (open requests only) |
| Start / complete trip | Assigned driver only |
| Cancel trip (after accept) | Owning customer or assigned driver (different cancel outcomes) |

---

## 2. TripRequest statuses

| Enum | Value | Terminal | Meaning |
|------|-------|----------|---------|
| `NO_DRIVER_FOUND` | 1 | No | Open in driver pool |
| `CUSTOMER_CANCELED` | 2 | Yes | Customer canceled before accept |
| `DRIVER_ACCEPTED` | 3 | Yes | Driver accepted — frozen forever |
| `EXPIRED` | 4 | Yes | No driver found in time (scheduler; not implemented) |

---

## 3. Trip statuses

| Enum | Value | Terminal | Meaning |
|------|-------|----------|---------|
| `TRIP_ACCEPTED` | 1 | No | Driver assigned, not started |
| `TRIP_IN_PROGRESS` | 2 | No | Journey started |
| `TRIP_COMPLETED` | 3 | Yes | Normal completion |
| `TRIP_CANCELLED_BY_CUSTOMER` | 4 | Yes | Customer canceled after accept |
| `TRIP_CANCELLED_BY_DRIVER` | 5 | Yes | Driver backed out after accept |

---

## 4. State transitions

### TripRequest (pre-accept only)

| From | To | Trigger | Actor | Endpoint |
|------|-----|---------|-------|----------|
| — | `NO_DRIVER_FOUND` | Create request | Customer | `POST /trip-requests` |
| `NO_DRIVER_FOUND` | `CUSTOMER_CANCELED` | Cancel | Customer | `DELETE /trip-requests/:trip_request_id` |
| `NO_DRIVER_FOUND` | `EXPIRED` | Timeout | System | *(not implemented)* |
| `NO_DRIVER_FOUND` | `DRIVER_ACCEPTED` | Accept (+ create trip) | Driver | `POST /trip-requests/:trip_request_id/accept` |

### Trip (post-accept only)

| From | To | Trigger | Actor | Endpoint |
|------|-----|---------|-------|----------|
| — | `TRIP_ACCEPTED` | Accept | Driver | `POST /trip-requests/:trip_request_id/accept` |
| `TRIP_ACCEPTED` | `TRIP_IN_PROGRESS` | Start journey | Driver | `POST /trips/:trip_id/start` |
| `TRIP_ACCEPTED` | `TRIP_CANCELLED_BY_CUSTOMER` | Cancel | Customer | `POST /trips/:trip_id/cancel` |
| `TRIP_ACCEPTED` | `TRIP_CANCELLED_BY_DRIVER` | Cancel | Driver | `POST /trips/:trip_id/cancel` |
| `TRIP_IN_PROGRESS` | `TRIP_COMPLETED` | Complete | Driver | `POST /trips/:trip_id/complete` |
| `TRIP_IN_PROGRESS` | `TRIP_CANCELLED_BY_CUSTOMER` | Cancel mid-ride | Customer | `POST /trips/:trip_id/cancel` |

---

## 5. HTTP API

| Method | Path | Auth | Handler | Status |
|--------|------|------|---------|--------|
| POST | `/customers/signup` | — | CustomerSignup | Done |
| POST | `/drivers/signup` | — | DriverSignup | Done |
| POST | `/trip-requests` | Customer JWT | RequestTrip | Done |
| GET | `/trip-requests/:trip_request_id` | Customer JWT | GetDetails | Done |
| DELETE | `/trip-requests/:trip_request_id` | Customer JWT | CancelTripRequest | Done (pre-accept only) |
| GET | `/trip-requests/open` | Driver JWT | ListOpenTripRequests | Done |
| POST | `/trip-requests/:trip_request_id/accept` | Driver JWT | AcceptTripRequest | Done |
| POST | `/trips/:trip_id/start` | Driver JWT | StartTrip | Done |
| POST | `/trips/:trip_id/complete` | Driver JWT | CompleteTrip | Done |
| POST | `/trips/:trip_id/cancel` | Customer or Driver JWT | CancelTrip | Done |

### Route ordering

Register `GET /trip-requests/open` **before** `GET /trip-requests/:trip_request_id`.

### Responses

- Create trip request: `201` `{ "trip_request": {...} }`
- Accept: `201` `{ "trip": {...}, "trip_request": {...} }`
- Start / complete / cancel trip: `200` `{ "trip": {...} }`
- Cancel trip request: `204`

---

## 6. Invariants

1. Open driver pool = `trip_request.status == NO_DRIVER_FOUND` only.
2. `trip_request` is **immutable after `DRIVER_ACCEPTED`** (status and fields).
3. One `trip` per `trip_request` (`trip_request_id` unique on `trips`).
4. Customer cancel before accept → `trip_request` only; no `trip` row.
5. Customer/driver cancel after accept → **`trip` only**; `trip_request` unchanged.
6. Driver cancel after accept → `trip` → `TRIP_CANCELLED_BY_DRIVER`; customer must create a **new** `trip_request` to retry.
7. Driver cancel allowed only from `TRIP_ACCEPTED` (before start).
8. Customer cancel after accept allowed from `TRIP_ACCEPTED` or `TRIP_IN_PROGRESS`.
9. On accept, `trip.customer_id` is set from `trip_request.customer_id` and never changes.

---

## 7. Implementation status

| Feature | Spec | Code | Tests |
|---------|------|------|-------|
| Customer signup | Yes | Done | E2E |
| Driver signup | Yes | Done | E2E |
| Request trip | Yes | Done | E2E |
| Cancel before accept | Yes | Done | E2E |
| List open (driver) | Yes | Done | E2E |
| Accept trip request | Yes | Done | Unit |
| Start trip | Yes | Done | Unit |
| Complete trip | Yes | Done | Unit |
| Cancel trip (customer/driver) | Yes | Done | Unit |
| Freeze trip_request on accept | Yes | Done | Unit |
| Expire trip request (scheduler) | Yes | Not started | — |

---

## 8. Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.1 | 2026-06-24 | Document entity fields, `customer_id` on trips, authorization requirements |
| 1.0 | 2026-06-23 | Initial spec: Option B routes, frozen trip_request on accept, trip lifecycle statuses |
