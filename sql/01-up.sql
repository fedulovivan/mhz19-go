PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 1
UPDATE condition_functions SET name='Nil' WHERE name='IsNil';

-- change 2
CREATE TABLE schema_version (
	version INTEGER
);

-- change 3
INSERT INTO schema_version VALUES(1);

COMMIT;