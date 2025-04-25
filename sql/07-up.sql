PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 7;

-- change 1
-- adding new system device
INSERT INTO devices (native_id, device_class_id, origin) VALUES('device-id-for-the-watcher-message', 7, 'migration-07');

-- change 2
-- adding new action
INSERT INTO action_functions VALUES(10, 'WatchChanges');

-- change 3
-- adding new condition
INSERT INTO condition_functions VALUES(13, 'LdmOlderThan');

COMMIT;