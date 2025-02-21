PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 6;

-- change 1
-- aligning name with interstellar/valves-manipulator
UPDATE device_classes SET name = 'valves-manipulator' WHERE name = 'valve-manipulator';

COMMIT;