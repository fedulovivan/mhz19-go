PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 1
DROP INDEX messages_device_id_index;

-- change 2
ALTER TABLE rules DROP COLUMN comments;

COMMIT;