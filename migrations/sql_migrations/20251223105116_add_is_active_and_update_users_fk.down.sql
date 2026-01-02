-- +migrate Down

BEGIN;

ALTER TABLE users
DROP CONSTRAINT fk_users_location;

ALTER TABLE users
ADD CONSTRAINT fk_users_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE users
DROP COLUMN is_active;

DROP INDEX IF EXISTS ux_users_phone_active;
DROP INDEX IF EXISTS ux_users_nik_active;
DROP INDEX IF EXISTS idx_users_is_active;

ALTER TABLE users
ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);

ALTER TABLE users
ADD CONSTRAINT users_nik_key UNIQUE (nik);

ALTER TABLE users
ADD CONSTRAINT unique_users_record
UNIQUE (location_id, name, phone_number, nik);

COMMIT;
