-- +migrate Up

BEGIN;

ALTER TABLE toddlers
ADD COLUMN deleted_at TIMESTAMPTZ,
ADD COLUMN deleted_by_id INT,
ADD COLUMN created_by_id INT,
ADD COLUMN updated_by_id INT;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_deleted_by
FOREIGN KEY (deleted_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_created_by
FOREIGN KEY (created_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_updated_by
FOREIGN KEY (updated_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE toddlers
DROP CONSTRAINT fk_toddlers_parent,
DROP CONSTRAINT fk_toddlers_location;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_parent
FOREIGN KEY (parent_id)
REFERENCES parents(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

ALTER TABLE toddlers
ADD CONSTRAINT fk_toddlers_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

ALTER TABLE toddlers
DROP CONSTRAINT unique_toddler_record;

CREATE UNIQUE INDEX ux_toddlers_active
ON toddlers (name, birthdate, parent_id, location_id)
WHERE deleted_at IS NULL;

COMMIT;