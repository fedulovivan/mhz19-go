PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 1
UPDATE condition_functions SET name='IsNil' WHERE name='Nil';

-- change 2, 3
DROP TABLE IF EXISTS schema_version;

COMMIT;