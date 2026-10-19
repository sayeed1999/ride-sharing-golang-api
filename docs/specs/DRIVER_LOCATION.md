# Driver Location Specification

**Version:** 1.0  
**Status:** Canonical — implementation source of truth  
**Base path (production):** `/api/trip`  
**Base path (e2e tests):** `/` (module mounted at root)

---

## 1. Domain model

| Entity | Purpose | Storage |
|--------|---------|---------|
| **DriverLocation** | Ephemeral live position of a driver | Redis hash + TTL |

No PostgreSQL table. Location is not persisted after TTL expiry.

### DriverLocation fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `lat` | float64 | Yes | -90 to 90 |
| `lng` | float64 | Yes | -180 to 180 |
| `heading` | float64 | No | 0–360 degrees |
| `speed` | float64 | No | m/s, ≥ 0 |
| `updated_at` | timestamptz | Yes | Set by API on write |

### Redis key

```
driver:loc:{driver_id}  →  HASH { lat, lng, heading?, speed?, updated_at }
TTL: 60s (refreshed on each write)
```

### Authorization

| Action | Who may perform it |
|--------|-------------------|
| Push location | Assigned driver with an active trip (`TRIP_ACCEPTED` or `TRIP_IN_PROGRESS`) |
| Read location | Owning customer or assigned driver on the same active trip |

---

## 2. HTTP API

| Method | Path | Auth | Handler | Status |
|--------|------|------|---------|--------|
| POST | `/drivers/location` | Driver JWT | UpdateDriverLocation | Done |
| GET | `/trips/:trip_id/driver-location` | Customer or Driver JWT | GetDriverLocation | Done |

### Request / response

**POST `/drivers/location`**

```json
{ "lat": 23.78, "lng": 90.41, "heading": 180, "speed": 12.5 }
```

- `204` on success
- `400` invalid coordinates or no active trip
- `401` / `403` auth failures

**GET `/trips/:trip_id/driver-location`**

```json
{ "driver_location": { "lat": 23.78, "lng": 90.41, "heading": 180, "speed": 12.5, "updated_at": "..." } }
```

- `200` on success
- `404` trip not found, location not available (TTL expired), or trip not trackable
- `403` actor not on trip

---

## 3. Invariants

1. Location writes only when driver has a trip in `TRIP_ACCEPTED` or `TRIP_IN_PROGRESS`.
2. Location reads only when trip is in `TRIP_ACCEPTED` or `TRIP_IN_PROGRESS`.
3. Driver may only write their own location (from JWT → driver middleware).
4. No location history in PostgreSQL.
5. Stale location (TTL expired) returns `404`.

---

## 4. Client contract

- Driver app: POST every 3–5s while on an active trip.
- Customer app: GET every 3–5s while trip is active.
- Stop pings when trip ends; TTL handles offline detection.

---

## 5. Implementation status

| Feature | Spec | Code | Tests |
|---------|------|------|-------|
| Push driver location | Yes | Done | Unit |
| Read driver location | Yes | Done | Unit + E2E |
| Redis storage + TTL | Yes | Done | E2E |

---

## 6. Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-06 | Initial spec: Redis ephemeral location, HTTP ping/poll |
