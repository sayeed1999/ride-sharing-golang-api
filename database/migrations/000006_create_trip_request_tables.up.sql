
CREATE TABLE IF NOT EXISTS "trip.trip_requests" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "customer_id" UUID NOT NULL,
    "origin" VARCHAR(255) NOT NULL,
    "destination" VARCHAR(255) NOT NULL,
    "status" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("customer_id") REFERENCES "trip.customers"("id")
);
