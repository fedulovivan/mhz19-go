PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 3;

-- change 1
CREATE INDEX raam_rule_id_index ON rule_action_argument_mappings(rule_id);

-- change 2
CREATE INDEX rcoaa_rule_id_index ON rule_condition_or_action_arguments(rule_id);

-- change 3
CREATE INDEX rc_rule_id_index ON rule_conditions(rule_id);

-- change 4
CREATE INDEX ra_rule_id_index ON rule_actions(rule_id);

COMMIT;