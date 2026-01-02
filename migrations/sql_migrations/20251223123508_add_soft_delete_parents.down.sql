-- +migrate Down

BEGIN;

DROP INDEX IF EXISTS ux_parents_active;

ALTER TABLE parents
DROP CONSTRAINT fk_parents_deleted_by,
DROP CONSTRAINT fk_parents_created_by,
DROP CONSTRAINT fk_parents_updated_by,
DROP CONSTRAINT fk_parents_location;

ALTER TABLE parents
DROP COLUMN deleted_at,
DROP COLUMN deleted_by_id,
DROP COLUMN created_by_id,
DROP COLUMN updated_by_id;

ALTER TABLE parents
ADD CONSTRAINT fk_parents_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

ALTER TABLE parents
ADD CONSTRAINT unique_parents_record
UNIQUE (location_id, name, phone_number, nik);

COMMIT;