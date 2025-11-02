ALTER TABLE trip.drivers DROP COLUMN auth_user_id;
ALTER TABLE trip.drivers ADD COLUMN auth_user_id UUID;
ALTER TABLE trip.drivers ADD CONSTRAINT fk_auth_user
FOREIGN KEY (auth_user_id) REFERENCES auth.users(id) ON DELETE CASCADE;