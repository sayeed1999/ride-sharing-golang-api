# PROJECT STRUCTURE

```bash
ride-sharing-api/
├── cmd/
│   ├── api/
│   │   └── main.go                    # Entry point for core API
│   ├── trip-processor/
│   │   └── main.go                    # Separate service for heavy workloads like matchmaking
│   └── migrate/
│       └── main.go                    # Separate tool for running database migrations
│
├── internal/                          # Private application code
│   │
│   ├── app/                           # Application modules (each = potential microservice)
│   │   │
│   │   ├── auth/                      # AUTH MODULE (Generic, reusable)
│   │   │   ├── domain/
│   │   │   │   ├── user.go
│   │   │   │   ├── role.go
│   │   │   │   ├── user_role.go
│   │   │   │   └── errors.go
│   │   │   ├── repository/
│   │   │   │   ├── repository.go      # Interface
│   │   │   │   └── postgres/
│   │   │   │       ├── user_repo.go
│   │   │   │       └── role_repo.go
│   │   │   ├── service/
│   │   │   │   └── auth_service.go    # PUBLIC API (exposed to other modules)
│   │   │   ├── usecase/               # Business logic
│   │   │   │   ├── register.go
│   │   │   │   ├── login.go
│   │   │   │   ├── assign_role.go
│   │   │   │   └── verify_token.go
│   │   │   ├── handler/
│   │   │   │   ├── http/
│   │   │   │   │   ├── handler.go     # HTTP handlers
│   │   │   │   │   ├── routes.go
│   │   │   │   │   └── dto.go
│   │   │   │   └── grpc/              # Future: gRPC interface
│   │   │   │       └── auth.proto
│   │   │   └── module.go              # Module initialization
│   │   │
│   │   ├── rider/                     # RIDER MODULE (Customer operations + Trip requests)
│   │   │   ├── domain/
│   │   │   │   ├── rider.go           # Rider entity
│   │   │   │   ├── trip.go            # Trip entity (rider perspective)
│   │   │   │   ├── trip_request.go    # Ride request
│   │   │   │   ├── cancellation.go    # Cancellation logic
│   │   │   │   ├── repository.go
│   │   │   │   └── errors.go
│   │   │   ├── repository/
│   │   │   │   └── postgres/
│   │   │   │       ├── rider_repo.go
│   │   │   │       ├── trip_repo.go
│   │   │   │       └── trip_request_repo.go
│   │   │   ├── service/
│   │   │   │   └── rider_service.go   # PUBLIC API
│   │   │   │       # Methods:
│   │   │   │       # - CreateProfile(authUserID, phone, address)
│   │   │   │       # - GetProfile(riderID)
│   │   │   │       # - RequestRide(riderID, pickup, dropoff) → TripRequest
│   │   │   │       # - CancelRide(riderID, tripID, reason)
│   │   │   │       # - GetRideHistory(riderID) → []Trip
│   │   │   │       # - GetActiveTrip(riderID) → Trip
│   │   │   │       # - GetTripDetails(riderID, tripID) → Trip
│   │   │   ├── usecase/
│   │   │   │   ├── create_profile.go
│   │   │   │   ├── request_ride.go    # US-C02 - Creates request, publishes event
│   │   │   │   ├── cancel_ride.go     # US-C03, US-C04
│   │   │   │   ├── get_ride_history.go # US-C05
│   │   │   │   └── update_trip_status.go # Handles trip updates from processing service
│   │   │   ├── handler/
│   │   │   │   ├── http/
│   │   │   │   │   ├── handler.go
│   │   │   │   │   ├── routes.go
│   │   │   │   │   └── dto.go
│   │   │   │   └── graphql/
│   │   │   │       ├── resolver.go
│   │   │   │       └── schema.graphql
│   │   │   └── module.go
│   │   │
│   │   ├── driver/                    # DRIVER MODULE (Driver operations + Trip fulfillment)
│   │   │   ├── domain/
│   │   │   │   ├── driver.go
│   │   │   │   ├── availability.go
│   │   │   │   ├── vehicle.go
│   │   │   │   ├── trip.go            # Trip entity (driver perspective)
│   │   │   │   ├── trip_assignment.go # Driver-trip assignment
│   │   │   │   ├── repository.go
│   │   │   │   └── errors.go
│   │   │   ├── repository/
│   │   │   │   └── postgres/
│   │   │   │       ├── driver_repo.go
│   │   │   │       ├── availability_repo.go
│   │   │   │       ├── trip_repo.go
│   │   │   │       └── trip_assignment_repo.go
│   │   │   ├── service/
│   │   │   │   └── driver_service.go  # PUBLIC API
│   │   │   │       # Methods:
│   │   │   │       # - CreateProfile(authUserID, license, vehicle)
│   │   │   │       # - GetProfile(driverID)
│   │   │   │       # - ToggleAvailability(driverID, status)
│   │   │   │       # - GetAvailableDrivers(location, radius) → []Driver
│   │   │   │       # - GetNearbyRequests(driverID, location) → []TripRequest
│   │   │   │       # - AcceptTrip(driverID, tripID) → Trip
│   │   │   │       # - StartTrip(driverID, tripID) → Trip
│   │   │   │       # - CompleteTrip(driverID, tripID) → Trip
│   │   │   │       # - GetRideHistory(driverID) → []Trip
│   │   │   │       # - GetActiveTrip(driverID) → Trip
│   │   │   ├── usecase/
│   │   │   │   ├── create_profile.go
│   │   │   │   ├── toggle_availability.go # US-D02
│   │   │   │   ├── get_nearby_requests.go # US-D03
│   │   │   │   ├── accept_trip.go     # US-D04 - Publishes acceptance event
│   │   │   │   ├── start_trip.go      # Trip started
│   │   │   │   ├── complete_trip.go   # US-D05
│   │   │   │   ├── get_ride_history.go # US-D06
│   │   │   │   └── update_trip_assignment.go # Handles assignments from processing service
│   │   │   ├── handler/
│   │   │   │   ├── http/
│   │   │   │   │   ├── handler.go
│   │   │   │   │   ├── routes.go
│   │   │   │   │   └── dto.go
│   │   │   │   └── graphql/
│   │   │   │       ├── resolver.go
│   │   │   │       └── schema.graphql
│   │   │   └── module.go
│   │   │
│   │   ├── tripprocessor/            # TRIP PROCESSING MODULE (Core algorithms - internal only)
│   │   │   ├── domain/
│   │   │   │   ├── trip_state.go      # Trip state machine
│   │   │   │   ├── trip_status.go     # Status transitions
│   │   │   │   ├── match_result.go
│   │   │   │   ├── driver_score.go
│   │   │   │   ├── repository.go
│   │   │   │   └── errors.go
│   │   │   ├── repository/
│   │   │   │   └── postgres/
│   │   │   │       ├── trip_state_repo.go
│   │   │   │       └── match_history_repo.go
│   │   │   ├── service/
│   │   │   │   └── trip_processing_service.go # INTERNAL API (not exposed via HTTP)
│   │   │   │       # Methods (called via events only):
│   │   │   │       # - ProcessTripRequest(requestID) → MatchResult
│   │   │   │       # - ValidateTripAcceptance(tripID, driverID) → bool
│   │   │   │       # - UpdateTripState(tripID, newStatus) → TripState
│   │   │   │       # - RecalculateMatch(tripID) → MatchResult
│   │   │   ├── usecase/
│   │   │   │   ├── process_trip_request.go
│   │   │   │   ├── validate_acceptance.go
│   │   │   │   ├── handle_cancellation.go
│   │   │   │   └── update_state.go
│   │   │   ├── matching/
│   │   │   │   ├── strategy.go        # Matching interface
│   │   │   │   ├── nearest_driver.go  # Distance-based matching
│   │   │   │   ├── best_rated.go      # Rating-based matching
│   │   │   │   ├── balanced.go        # Hybrid strategy
│   │   │   │   └── scorer.go          # Driver scoring algorithm
│   │   │   ├── statemachine/
│   │   │   │   ├── transitions.go     # Valid state transitions
│   │   │   │   └── validator.go       # State validation
│   │   │   ├── scheduler/
│   │   │   │   ├── retry_matching.go  # Retry failed matches
│   │   │   │   └── timeout_handler.go # Handle request timeouts
│   │   │   ├── handler/
│   │   │   │   └── event/             # NO HTTP/GraphQL - event-driven only
│   │   │   │       ├── subscriber.go  # Listens to TripRequested, TripAccepted, etc.
│   │   │   │       └── publisher.go   # Publishes MatchFound, TripStateChanged
│   │   │   └── module.go
│   │   │
│   │   ├── payment/                   # PAYMENT MODULE (Future)
│   │   │   ├── domain/
│   │   │   │   ├── payment.go
│   │   │   │   ├── transaction.go
│   │   │   │   └── repository.go
│   │   │   ├── repository/
│   │   │   │   └── postgres/
│   │   │   │       └── payment_repo.go
│   │   │   ├── service/
│   │   │   │   └── payment_service.go # PUBLIC API
│   │   │   │       # Methods:
│   │   │   │       # - ProcessPayment(tripID, amount) → Payment
│   │   │   │       # - RefundPayment(paymentID) → Payment
│   │   │   │       # - GetPaymentStatus(paymentID) → Payment
│   │   │   ├── usecase/
│   │   │   │   ├── process_payment.go
│   │   │   │   └── refund_payment.go
│   │   │   ├── gateway/
│   │   │   │   ├── stripe.go
│   │   │   │   └── gateway.go         # Interface
│   │   │   ├── handler/
│   │   │   │   ├── http/
│   │   │   │   │   └── handler.go
│   │   │   │   └── webhook/
│   │   │   │       └── stripe_webhook.go
│   │   │   └── module.go
│   │   │
│   │   ├── pricing/                   # PRICING MODULE (Future)
│   │   │   ├── domain/
│   │   │   │   ├── fare.go
│   │   │   │   ├── pricing_strategy.go
│   │   │   │   └── vehicle_pricing.go
│   │   │   ├── repository/
│   │   │   │   └── postgres/
│   │   │   │       └── pricing_config_repo.go
│   │   │   ├── service/
│   │   │   │   └── pricing_service.go # PUBLIC API
│   │   │   │       # Methods:
│   │   │   │       # - CalculateFare(distance, duration, vehicleType) → Fare
│   │   │   │       # - EstimateFare(pickup, dropoff, vehicleType) → Fare
│   │   │   ├── usecase/
│   │   │   │   ├── calculate_fare.go
│   │   │   │   └── estimate_fare.go
│   │   │   └── module.go
│   │   │
│   │   └── notification/              # NOTIFICATION MODULE (Future)
│   │       ├── domain/
│   │       │   ├── notification.go
│   │       │   └── repository.go
│   │       ├── repository/
│   │       │   └── postgres/
│   │       │       └── notification_repo.go
│   │       ├── service/
│   │       │   └── notification_service.go # PUBLIC API
│   │       │       # Methods:
│   │       │       # - SendSMS(phone, message)
│   │       │       # - SendPush(userID, title, body)
│   │       │       # - SendEmail(email, subject, body)
│   │       ├── usecase/
│   │       │   └── send_notification.go
│   │       ├── provider/
│   │       │   ├── twilio.go
│   │       │   └── firebase.go
│   │       └── module.go
│   │
│   └── shared/                        # Internal shared code across modules
│       ├── event/
│       │   ├── bus.go                 # Event bus interface
│       │   ├── types.go               # Event definitions:
│       │   │                          # - TripRequested (from rider)
│       │   │                          # - MatchFound (from trip-processing)
│       │   │                          # - TripAccepted (from driver)
│       │   │                          # - TripStarted (from driver)
│       │   │                          # - TripCompleted (from driver)
│       │   │                          # - TripCancelled (from rider/driver)
│       │   │                          # - TripStateChanged (from trip-processing)
│       │   ├── publisher.go
│       │   └── subscriber.go
│       ├── errors/
│       │   ├── errors.go
│       │   └── codes.go
│       └── types/
│           ├── location.go            # Shared Location type
│           ├── money.go               # Money value object
│           └── pagination.go
├── config/
│   └── config.go
├── database/
│   ├── postgres.go
│   └── transaction.go
├── logger/
│   └── logger.go
├── pkg/                               # Public shared utilities (can be imported externally)
│   ├── jwt/
│   │   └── jwt.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── httpclient/
│   │   └── client.go
│   └── validator/
│       └── validator.go
│
├── infrastructure/                    # Infrastructure setup
│   ├── http/
│   │   ├── server.go                  # HTTP server
│   │   └── router.go                  # Combines rider, driver, auth routes
│   ├── graphql/
│   │   ├── server.go
│   │   └── schema.go
│   ├── database/
│   │   └── migrations/
│   │       ├── 000001_auth_tables.up.sql
│   │       ├── 000002_rider_tables.up.sql
│   │       ├── 000003_driver_tables.up.sql
│   │       ├── 000004_trip_processing_tables.up.sql
│   │       └── ...
│   ├── cache/
│   │   └── redis.go
│   └── messaging/
│       ├── eventbus.go                # In-memory or Redis pub/sub
│       └── kafka.go                   # Future: Kafka for microservices
│
├── api/
│   ├── rest/
│   │   └── openapi.yaml
│   └── graphql/
│       └── schema.graphql
│
├── scripts/
│   ├── build.sh
│   ├── test.sh
│   ├── seed-roles.sh
│   └── migrate.sh
│
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── kubernetes/                    # For future microservices
│       ├── auth-service.yaml
│       ├── rider-service.yaml
│       ├── driver-service.yaml
│       └── trip-processing-service.yaml
│
├── docs/
│   ├── architecture.md
│   ├── module-contracts.md            # Module interface documentation
│   ├── event-flow.md                  # Event-driven flow diagrams
│   ├── trip-state-machine.md          # Trip status transitions
│   ├── migration-to-microservices.md
│   └── api-design.md
│
├── tests/
│   ├── integration/
│   │   ├── auth_module_test.go
│   │   ├── rider_module_test.go
│   │   ├── driver_module_test.go
│   │   └── trip_processing_test.go
│   └── e2e/
│       └── full_trip_flow_test.go
│
├── go.mod
├── go.sum
├── .env.example
├── Makefile
└── README.md
```

# KEY ARCHITECTURAL DECISIONS:

## 1. Trip Data Ownership
- **Rider Service**: Owns rider perspective of trips (trip requests, cancellations, history)
- **Driver Service**: Owns driver perspective of trips (assignments, acceptances, completions)
- **Trip-Processing Service**: Owns core trip state and orchestration logic

## 2. Communication Flow
```
Rider requests ride → TripRequested event → Trip-Processing
Trip-Processing matches driver → MatchFound event → Driver Service
Driver accepts → TripAccepted event → Trip-Processing → Update Rider
Driver starts → TripStarted event → Trip-Processing → Update Rider
Driver completes → TripCompleted event → Trip-Processing → Update both
```

## 3. No Direct Calls
- Rider service NEVER calls driver service directly
- Driver service NEVER calls rider service directly
- Trip-processing service has NO HTTP endpoints
- All coordination happens through events

## 4. Future Microservice Split
When splitting into microservices:
- Each service becomes independent deployment
- Event bus becomes Kafka/RabbitMQ
- Each service has its own database
- Trip-processing remains internal orchestrator
