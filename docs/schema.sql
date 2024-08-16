-- PRAGMA foreign_keys=OFF;
-- BEGIN TRANSACTION;

DROP TABLE IF EXISTS device_classes;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS rules;
DROP TABLE IF EXISTS rule_conditions;
DROP TABLE IF EXISTS rule_actions;
DROP TABLE IF EXISTS rule_condition_or_action_arguments;
DROP TABLE IF EXISTS rule_action_argument_mappings;

CREATE TABLE device_classes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE condition_functions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE action_functions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE devices (
   	id INTEGER PRIMARY KEY AUTOINCREMENT,
    native_id TEXT NOT NULL,
    device_class_id INTEGER NOT NULL,
	name TEXT,
	comments TEXT,
    CONSTRAINT devices_fk_dc FOREIGN KEY (device_class_id) REFERENCES device_classes(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	comments TEXT NOT NULL,
	is_disabled INTEGER,
	throttle INTEGER
);

CREATE TABLE rule_conditions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_type INTEGER,
	logic_or INTEGER,
	parent_condition_id INTEGER,
	CONSTRAINT rule_conditions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_conditions_fk_function FOREIGN KEY (function_type) REFERENCES condition_functions(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT rule_conditions_fk_parent FOREIGN KEY (parent_condition_id) REFERENCES rule_conditions(id)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_actions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_type INTEGER,
    device_id INTEGER,
    CONSTRAINT rule_actions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_actions_fk_function FOREIGN KEY (function_type) REFERENCES action_functions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_actions_fk_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_condition_or_action_arguments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    condition_id INTEGER,
    action_id INTEGER,
    argument_name TEXT NOT NULL,
    is_list INTEGER,
    value TEXT,
    device_id INTEGER,
    device_class_id INTEGER,
    CONSTRAINT rule_ca_arguments_fk_conditions FOREIGN KEY (condition_id) REFERENCES rule_conditions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_actions FOREIGN KEY (action_id) REFERENCES rule_actions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_dc FOREIGN KEY (device_class_id) REFERENCES device_classes(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_action_argument_mappings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    argument_id INTEGER,
    key TEXT,
    value TEXT,
    CONSTRAINT rule_action_mappings_fk_arguments FOREIGN KEY (argument_id) REFERENCES rule_condition_or_action_arguments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

--------------

-- migration 00
-- https://www.prisma.io/dataguide/sqlite/importing-and-exporting-data-in-sqlite
-- INSERT INTO sqlite_sequence VALUES('rules',1);
-- INSERT INTO sqlite_sequence VALUES('rule_conditions',1);
-- INSERT INTO sqlite_sequence VALUES('rule_condition_or_action_arguments',2);
-- INSERT INTO sqlite_sequence VALUES('rule_actions',1);

INSERT INTO device_classes VALUES(0,'<unknown>');
INSERT INTO device_classes VALUES(1,'zigbee-device');
INSERT INTO device_classes VALUES(2,'device-pinger');
INSERT INTO device_classes VALUES(3,'valve-manipulator');
INSERT INTO device_classes VALUES(4,'zigbee-bridge');
INSERT INTO device_classes VALUES(5,'telegram-bot');

INSERT INTO action_functions VALUES(1,'PostSonoffSwitchMessage');
INSERT INTO action_functions VALUES(2,'TelegramBotMessage');
INSERT INTO action_functions VALUES(3,'ValveSetState');
INSERT INTO action_functions VALUES(4,'YeelightDeviceSetPower');
INSERT INTO action_functions VALUES(5,'Zigbee2MqttSetState');

INSERT INTO condition_functions VALUES(1,'Changed');
INSERT INTO condition_functions VALUES(2,'Equal');
INSERT INTO condition_functions VALUES(3,'InList');
INSERT INTO condition_functions VALUES(4,'NotEqual');
INSERT INTO condition_functions VALUES(5,'NotNil');
INSERT INTO condition_functions VALUES(6,'ZigbeeDevice');

INSERT INTO rules VALUES(1,'test mapping 1',0,NULL);
INSERT INTO rule_conditions VALUES(1,1,2,NULL,NULL);
INSERT INTO rule_actions VALUES(1,1,2,NULL);
INSERT INTO rule_condition_or_action_arguments VALUES(1,1,NULL,'Left',NULL,'$deviceClass',NULL,NULL);
INSERT INTO rule_condition_or_action_arguments VALUES(2,1,NULL,'Right',NULL,NULL,NULL,2);

-- COMMIT;
-- migration 01 (sqlite does not support ADD CONSTRAINT)
-- https://www.alphacodingskills.com/sqlite/notes/sqlite-keyword-add-constraint.php
-- ALTER TABLE rule_conditions ADD CONSTRAINT function_name_check CHECK (function_name IN ('Changed','Equal','InList','NotEqual','NotNil','ZigbeeDevice'));
-- ALTER TABLE rule_actions ADD CONSTRAINT function_name_check CHECK (function_name IN ('PostSonoffSwitchMessage','YeelightDeviceSetPower','Zigbee2MqttSetState','ValveSetState','TelegramBotMessage'));