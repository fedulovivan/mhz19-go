PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 8;

-- change 1
ALTER TABLE rule_conditions ADD COLUMN is_disabled INTEGER NOT NULL DEFAULT 0;

-- change 2
ALTER TABLE rule_actions ADD COLUMN is_disabled INTEGER NOT NULL DEFAULT 0;

COMMIT;