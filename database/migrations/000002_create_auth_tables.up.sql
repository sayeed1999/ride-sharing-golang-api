CREATE TABLE "auth"."roles" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL
);

CREATE TABLE "auth"."users" (
  "id" SERIAL PRIMARY KEY,
  "email" varchar(255) UNIQUE NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "password_salt" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "auth"."user_roles" (
  "id" SERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  UNIQUE ("user_id", "role_id"),
  FOREIGN KEY ("user_id") REFERENCES "auth"."users" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "auth"."roles" ("id")
);
