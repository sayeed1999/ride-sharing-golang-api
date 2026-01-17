CREATE TABLE "trip"."customers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" varchar(255) UNIQUE NOT NULL,
  "name" varchar(255),
  "auth_user_id" int,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trip"."drivers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" varchar(255) UNIQUE NOT NULL,
  "name" varchar(255),
  "auth_user_id" int,
  "vehicle_type_enum_code" int,
  "vehicle_registration" varchar(100),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "trip"."vehicle_types" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(50) UNIQUE NOT NULL,
  "enum_code" int UNIQUE NOT NULL
);
