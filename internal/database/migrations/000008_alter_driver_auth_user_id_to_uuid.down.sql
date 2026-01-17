ALTER TABLE trip.drivers DROP CONSTRAINT IF EXISTS fk_auth_user;
ALTER TABLE trip.drivers DROP COLUMN auth_user_id;
ALTER TABLE trip.drivers ADD COLUMN auth_user_id INTEGER;