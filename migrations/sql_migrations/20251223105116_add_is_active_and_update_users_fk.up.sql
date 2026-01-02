-- +migrate Up

BEGIN;

ALTER TABLE users
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE users
DROP CONSTRAINT fk_users_location;

ALTER TABLE users
ADD CONSTRAINT fk_users_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON DELETE RESTRICT
ON UPDATE CASCADE;

ALTER TABLE users
DROP CONSTRAINT users_phone_number_key;

ALTER TABLE users
DROP CONSTRAINT users_nik_key;

ALTER TABLE users
DROP CONSTRAINT unique_users_record;

CREATE UNIQUE INDEX ux_users_phone_active
ON users (phone_number)
WHERE is_active = TRUE;

CREATE UNIQUE INDEX ux_users_nik_active
ON users (nik)
WHERE is_active = TRUE;

COMMIT;
