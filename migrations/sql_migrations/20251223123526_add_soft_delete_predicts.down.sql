-- +migrate Down

BEGIN;

ALTER TABLE predicts
DROP CONSTRAINT fk_predicts_deleted_by,
DROP CONSTRAINT fk_predicts_created_by,
DROP CONSTRAINT fk_predicts_toddler,
DROP CONSTRAINT fk_predicts_location;

ALTER TABLE predicts
DROP COLUMN deleted_at,
DROP COLUMN deleted_by_id,
DROP COLUMN created_by_id;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_toddler
FOREIGN KEY (toddler_id)
REFERENCES toddlers(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

ALTER TABLE predicts
ADD CONSTRAINT fk_predicts_location
FOREIGN KEY (location_id)
REFERENCES locations(id)
ON UPDATE CASCADE
ON DELETE CASCADE;

COMMIT;
