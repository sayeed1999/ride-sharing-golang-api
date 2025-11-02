# Contributing Guidelines

Thank you for your interest in contributing to our project! To ensure a smooth and collaborative process, please follow these guidelines when adding new features, fixing bugs, or making any other changes.

## Adding a New Endpoint

This guide outlines the steps to add a new API endpoint to the application, following our modular monolith architecture. Each step builds upon the previous one, ensuring a structured approach.

### 1. Database Migration

- If your endpoint requires database changes, create a new migration file in `database/migrations/`.
- **Naming Convention:** Use an incremental number prefix followed by a descriptive name (e.g., `000007_create_new_feature_table.up.sql`). While `YYYYMMDDHHMMSS` is a common alternative, we currently use incremental numbering.
- Write the SQL statements for creating or altering tables in the `.up.sql` file and the statements to revert those changes in the `.down.sql` file.
- Run `scripts/migrate.sh` to apply the migration. Always verify the migration runs successfully.

### 2. Domain Model Definition

- Define the core data structures and business entities for your new feature in the appropriate module under `internal/app/<module>/domain/`.
- These models should be pure Go structs, representing the business domain, and should be independent of persistence or transport concerns.

### 3. Repository Layer

- **Interface (`internal/app/<module>/repository/`):** Define an interface that declares the methods for interacting with the persistence layer (e.g., `SaveUser(user *domain.User) error`, `GetUserByID(id string) (*domain.User, error)`).
- **Implementation (`internal/app/<module>/repository/postgres/`):** Implement the defined interface using PostgreSQL-specific logic. This is where you'll write SQL queries or use an ORM/query builder. Ensure proper error handling and transaction management.

### 4. Usecase Layer (Business Logic)

- Implement the business logic and application-specific rules in `internal/app/<module>/usecase/`.
- A usecase should represent a single, cohesive business operation (e.g., `RegisterUser`, `CreateTrip`).
- It orchestrates interactions between domain models and repositories. Usecases should be independent of HTTP details.
- **Functional Testing:** Usecases should be thoroughly functional tested to ensure the correctness of the business logic, typically by mocking the repository interfaces.

### 5. API Exposure Layer

This layer is responsible for exposing the application's functionality through various API protocols. Depending on the requirements, you might implement one or more of the following:

#### 5.1. HTTP Endpoints

- **DTOs (`internal/app/<module>/handler/http/dto.go`):** Define Data Transfer Objects for request and response payloads. These DTOs handle data serialization/deserialization and validation for the API.
- **Handler (`internal/app/<module>/handler/http/`):** Create an HTTP handler function or method that:
    - Parses and validates the incoming HTTP request using the defined DTOs.
    - Calls the appropriate usecase to execute the business logic.
    - Translates the usecase's result into an HTTP response, including proper status codes and error messages.
- **Routes (`internal/app/<module>/handler/http/routes.go`):** Register the new API endpoint by defining its path, HTTP method (GET, POST, PUT, DELETE), and associating it with your handler function.

#### 5.2. GraphQL Endpoints

- **Schema Definition (`internal/app/<module>/handler/graphql/schema.graphql`):** Define your GraphQL types, queries, mutations, and subscriptions in a `.graphql` schema file.
- **Resolvers (`internal/app/<module>/handler/graphql/resolver.go`):** Implement resolver functions that fetch data for the fields defined in your schema. Resolvers should interact with the usecase layer to retrieve and manipulate data.
- **Module Integration:** Ensure your GraphQL schema and resolvers are correctly integrated into the main GraphQL server setup, typically within `internal/app/<module>/module.go` or a dedicated GraphQL module.

#### 5.3. gRPC Endpoints

- **Protocol Buffer Definition (`internal/app/<module>/handler/grpc/<service>.proto`):** Define your gRPC services and messages using Protocol Buffers. This file specifies the service interface and the structure of the request and response messages.
- **Service Implementation (`internal/app/<module>/handler/grpc/service.go`):** Implement the gRPC service interface generated from your `.proto` file. This service should call the appropriate usecase to perform business logic.
- **Server Registration:** Register your gRPC service with the main gRPC server, typically within `internal/app/<module>/module.go` or a dedicated gRPC module.

### 6. Module Definition and Dependency Injection

- In `internal/app/<module>/module.go`, ensure that all dependencies (repositories, usecases, handlers) are correctly wired together using dependency injection.
- This file is also responsible for exposing the module's HTTP routes to the main application router.

### 7. Testing

Testing is crucial for maintaining code quality and ensuring the reliability of our application. We categorize our tests into two main types:

#### 7.1. Unit Testing

- Unit tests focus on individual components (functions, methods, structs) in isolation.
- They should be placed alongside the code they test (e.g., `usecase/login_test.go` for `usecase/login.go`).
- Aim for high code coverage and ensure all critical logic paths are tested.
- Use mocks or stubs for external dependencies (e.g., database, external services) to keep unit tests fast and focused.

#### 7.2. End-to-End Testing

- End-to-end tests verify the entire flow of a feature, from the API request to database interactions and the final response.
- These tests are located in the `tests/e2e/` directory.
- Create a new test file named `<feature>_test.go` (e.g., `user_registration_test.go`).
- End-to-end tests should simulate client requests and assert the correctness of the API's behavior under various scenarios.
