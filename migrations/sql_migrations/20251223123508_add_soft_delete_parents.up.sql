-- +migrate Up

BEGIN;

ALTER TABLE parents
ADD COLUMN deleted_at TIMESTAMPTZ,
ADD COLUMN deleted_by_id INT,
ADD COLUMN created_by_id INT,
ADD COLUMN updated_by_id INT;

ALTER TABLE parents
ADD CONSTRAINT fk_parents_deleted_by
FOREIGN KEY (deleted_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE parents
ADD CONSTRAINT fk_parents_created_by
FOREIGN KEY (created_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE parents
ADD CONSTRAINT fk_parents_updated_by
FOREIGN KEY (updated_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE parents
DROP CONSTRAINT fk_parents_location;

ALTER TABLE parents
ADD CONSTRAINT fk_parents_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

ALTER TABLE parents
DROP CONSTRAINT unique_parents_record;

CREATE UNIQUE INDEX ux_parents_active
ON parents (location_id, phone_number, nik)
WHERE deleted_at IS NULL;

COMMIT;