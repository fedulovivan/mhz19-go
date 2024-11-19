PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

-- undo change 3
UPDATE action_functions SET name = 'Zigbee2MqttSetState' WHERE id = 5;

-- undo change 2
INSERT INTO action_functions VALUES(3, 'ValveSetState');

-- undo change 1
UPDATE rule_actions SET function_type = 3 WHERE function_type = 5 AND rule_id in (SELECT id FROM rules WHERE name LIKE '%valve%');

-- undo change 0
UPDATE schema_version SET version = 3;

COMMIT;