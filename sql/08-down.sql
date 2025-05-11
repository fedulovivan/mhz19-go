PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- undo change 2
ALTER TABLE rule_actions DROP COLUMN is_disabled;

-- undo change 1
ALTER TABLE rule_conditions DROP COLUMN is_disabled;

-- undo change 0
UPDATE schema_version SET version = 7;

COMMIT;