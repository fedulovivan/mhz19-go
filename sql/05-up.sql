PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 5;

-- change 1
-- fixing bug from 00-up.sql
UPDATE device_classes SET id = 8 WHERE name = 'sonoff-announce';

-- change 2
-- insert new class
INSERT INTO device_classes VALUES(9, 'espresence-device');

COMMIT;