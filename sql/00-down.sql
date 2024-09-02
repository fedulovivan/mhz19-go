
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
DROP TABLE IF EXISTS device_classes;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS channel_types;
DROP TABLE IF EXISTS condition_functions;
DROP TABLE IF EXISTS action_functions;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS rules;
DROP TABLE IF EXISTS rule_conditions;
DROP TABLE IF EXISTS rule_actions;
DROP TABLE IF EXISTS rule_condition_or_action_arguments;
DROP TABLE IF EXISTS rule_action_argument_mappings;
COMMIT;
VACUUM;