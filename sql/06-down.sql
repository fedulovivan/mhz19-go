PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- undo change 1
UPDATE device_classes SET name = 'valve-manipulator' WHERE name = 'valves-manipulator';

-- undo change 0
UPDATE schema_version SET version = 5;

COMMIT;