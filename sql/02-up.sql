PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 1
CREATE INDEX messages_device_id_index ON messages (device_id);

-- change 2
ALTER TABLE rules ADD COLUMN comments TEXT;

COMMIT;