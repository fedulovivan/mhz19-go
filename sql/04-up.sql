PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- change 0
UPDATE schema_version SET version = 4;

-- change 1
-- replace ValveSetState with new MqttSetState for all actions
UPDATE rule_actions SET function_type = 5 WHERE function_type = 3;

-- change 2
-- delete ValveSetState acction
DELETE FROM action_functions WHERE id = 3;

-- change 3
-- rename Zigbee2MqttSetState to MqttSetState
UPDATE action_functions SET name = 'MqttSetState' WHERE id = 5;

COMMIT;