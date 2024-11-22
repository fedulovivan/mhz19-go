PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- undo change 2
DELETE FROM device_classes WHERE id = 9;

-- undo change 1
UPDATE device_classes SET id = 9 WHERE name = 'sonoff-announce';

-- undo change 0
UPDATE schema_version SET version = 4;

COMMIT;