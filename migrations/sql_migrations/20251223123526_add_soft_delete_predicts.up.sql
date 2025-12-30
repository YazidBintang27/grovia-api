-- +migrate Up

BEGIN;

ALTER TABLE predicts
ADD COLUMN deleted_at TIMESTAMPTZ,
ADD COLUMN deleted_by_id INT,
ADD COLUMN created_by_id INT;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_deleted_by
FOREIGN KEY (deleted_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_created_by
FOREIGN KEY (created_by_id)
REFERENCES users(id)
ON DELETE RESTRICT;

ALTER TABLE predicts
DROP CONSTRAINT fk_predicts_toddler,
DROP CONSTRAINT fk_predicts_location;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_toddler
FOREIGN KEY (toddler_id)
REFERENCES toddlers(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE RESTRICT;

COMMIT;
