PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 2;

-- change 1
DROP INDEX raam_rule_id_index;

-- change 2
DROP INDEX rcoaa_rule_id_index;

-- change 3
DROP INDEX rc_rule_id_index;

-- change 4
DROP INDEX ra_rule_id_index;

COMMIT;