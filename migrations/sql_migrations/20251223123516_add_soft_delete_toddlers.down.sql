-- +migrate Up

BEGIN;

DROP INDEX IF EXISTS ux_toddlers_active;

ALTER TABLE toddlers
DROP CONSTRAINT fk_toddlers_deleted_by,
DROP CONSTRAINT fk_toddlers_created_by,
DROP CONSTRAINT fk_toddlers_updated_by,
DROP CONSTRAINT fk_toddlers_parent,
DROP CONSTRAINT fk_toddlers_location;

ALTER TABLE toddlers
DROP COLUMN deleted_at,
DROP COLUMN deleted_by_id,
DROP COLUMN created_by_id,
DROP COLUMN updated_by_id;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_parent
FOREIGN KEY (parent_id)
REFERENCES parents(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

ALTER TABLE toddlers
ADD CONSTRAINT unique_toddler_record
UNIQUE (name, birthdate, parent_id, location_id);

COMMIT;