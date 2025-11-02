ALTER TABLE trip.customers DROP CONSTRAINT IF EXISTS fk_auth_user;
ALTER TABLE trip.customers DROP COLUMN auth_user_id;
ALTER TABLE trip.customers ADD COLUMN auth_user_id INTEGER;