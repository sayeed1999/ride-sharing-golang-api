CREATE TABLE IF NOT EXISTS "trip"."trips" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "trip_request_id" UUID NOT NULL UNIQUE,
    "customer_id" UUID NOT NULL,
    "driver_id" UUID NOT NULL,
    "status" INTEGER NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("trip_request_id") REFERENCES "trip"."trip_requests"("id"),
    FOREIGN KEY ("customer_id") REFERENCES "trip"."customers"("id"),
    FOREIGN KEY ("driver_id") REFERENCES "trip"."drivers"("id")
);
