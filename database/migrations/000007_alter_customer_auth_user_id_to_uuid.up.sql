ALTER TABLE trip.customers DROP COLUMN auth_user_id;
ALTER TABLE trip.customers ADD COLUMN auth_user_id UUID;
ALTER TABLE trip.customers ADD CONSTRAINT fk_auth_user
FOREIGN KEY (auth_user_id) REFERENCES auth.users(id) ON DELETE CASCADE;