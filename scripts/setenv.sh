export Server__Host=0.0.0.0
export Server__Port=8080

export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export POSTGRES_HOST=0.0.0.0 # use 'db' if using Docker Compose
export POSTGRES_PORT=5432
export POSTGRES_DB=ride_sharing_db

export PGADMIN_PASSWORD=admin
export PGADMIN_EMAIL=admin@local.com

export REQUIRE_ROLE_ON_REGISTRATION=true

# JWT secret used to sign access tokens in local dev. Replace in production.
export JWT_SECRET=dev_jwt_secret_change_me
