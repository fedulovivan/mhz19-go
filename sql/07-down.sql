PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- undo change 3
DELETE FROM condition_functions WHERE id = 13;

-- undo change 2
DELETE FROM action_functions WHERE id = 10;

-- undo change 1
DELETE FROM devices WHERE native_id = 'device-id-for-the-watcher-message';

-- undo change 0
UPDATE schema_version SET version = 6;

COMMIT;