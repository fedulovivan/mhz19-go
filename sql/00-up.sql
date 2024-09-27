PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

CREATE TABLE device_classes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE channel_types (
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
    native_id TEXT NOT NULL UNIQUE,
    device_class_id INTEGER NOT NULL,
	name TEXT,
	comments TEXT,
    origin TEXT,
    json TEXT,
    buried_timeout INTEGER,
    CONSTRAINT devices_fk_dc FOREIGN KEY (device_class_id) REFERENCES device_classes(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	is_disabled INTEGER,
	throttle_ms INTEGER
);

CREATE TABLE rule_conditions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_type INTEGER,
	logic_or INTEGER,
    function_inverted INTEGER,
	parent_condition_id INTEGER,
    other_device_id TEXT,
	CONSTRAINT rule_conditions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_conditions_fk_function FOREIGN KEY (function_type) REFERENCES condition_functions(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT rule_conditions_fk_parent FOREIGN KEY (parent_condition_id) REFERENCES rule_conditions(id)  ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_conditions_fk_devices FOREIGN KEY (other_device_id) REFERENCES devices(native_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_actions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_type INTEGER,
    CONSTRAINT rule_actions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_actions_fk_function FOREIGN KEY (function_type) REFERENCES action_functions(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_condition_or_action_arguments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
    condition_id INTEGER,
    action_id INTEGER,
    argument_name TEXT NOT NULL,
    is_list INTEGER,
    value TEXT,
    value_data_type TEXT,
    device_id TEXT,
    device_class_id INTEGER,
    channel_type_id INTEGER,
    CONSTRAINT rule_ca_arguments_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_conditions FOREIGN KEY (condition_id) REFERENCES rule_conditions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_actions FOREIGN KEY (action_id) REFERENCES rule_actions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_devices FOREIGN KEY (device_id) REFERENCES devices(native_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_dc FOREIGN KEY (device_class_id) REFERENCES device_classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_channels FOREIGN KEY (channel_type_id) REFERENCES channel_types(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_action_argument_mappings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
    argument_id INTEGER,
    key TEXT,
    value TEXT,
    CONSTRAINT rule_actions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_action_mappings_fk_arguments FOREIGN KEY (argument_id) REFERENCES rule_condition_or_action_arguments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE messages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_type_id INTEGER,
	device_class_id INTEGER,
    device_id TEXT,
	timestamp DATETIME NOT NULL,
	json TEXT NOT NULL,
    CONSTRAINT messages_fk_devices FOREIGN KEY (device_id) REFERENCES devices(native_id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT messages_fk_dc FOREIGN KEY (device_class_id) REFERENCES device_classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT messages_fk_channels FOREIGN KEY (channel_type_id) REFERENCES channel_types(id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO channel_types VALUES(1,'mqtt');
INSERT INTO channel_types VALUES(2,'telegram');
INSERT INTO channel_types VALUES(3,'dns-sd');
INSERT INTO channel_types VALUES(4,'system');
INSERT INTO channel_types VALUES(5,'rest');

INSERT INTO device_classes VALUES(1,'zigbee-device');
INSERT INTO device_classes VALUES(2,'device-pinger');
INSERT INTO device_classes VALUES(3,'valve-manipulator');
INSERT INTO device_classes VALUES(4,'zigbee-bridge');
INSERT INTO device_classes VALUES(5,'telegram-bot');
INSERT INTO device_classes VALUES(6,'sonoff-diy-plug');
INSERT INTO device_classes VALUES(7,'system');
INSERT INTO device_classes VALUES(9,'sonoff-announce');

INSERT INTO action_functions VALUES(1,'PostSonoffSwitchMessage');
INSERT INTO action_functions VALUES(2,'TelegramBotMessage');
INSERT INTO action_functions VALUES(3,'ValveSetState');
INSERT INTO action_functions VALUES(4,'YeelightDeviceSetPower');
INSERT INTO action_functions VALUES(5,'Zigbee2MqttSetState');
INSERT INTO action_functions VALUES(6,'RecordMessage');
INSERT INTO action_functions VALUES(7,'UpsertZigbeeDevices');
INSERT INTO action_functions VALUES(8,'UpsertSonoffDevice');
INSERT INTO action_functions VALUES(9,'PlayAlert');

INSERT INTO condition_functions VALUES(1,'Changed');
INSERT INTO condition_functions VALUES(2,'Equal');
INSERT INTO condition_functions VALUES(3,'InList');
INSERT INTO condition_functions VALUES(5,'IsNil');
INSERT INTO condition_functions VALUES(6,'ZigbeeDevice');
INSERT INTO condition_functions VALUES(7,'DeviceClass');
INSERT INTO condition_functions VALUES(8,'Channel');
INSERT INTO condition_functions VALUES(9,'FromEndDevice');
INSERT INTO condition_functions VALUES(10,'True');
INSERT INTO condition_functions VALUES(11,'False');
INSERT INTO condition_functions VALUES(12,'DeviceId');

INSERT INTO devices VALUES(1, '192.168.88.1', 2, 'MIKROTIK_ROUTER', NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(2, '192.168.88.44', 2, 'IPHONE_15_PRO_IP', NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(3, '192.168.0.11', 2, 'IPHONE_15_PRO_AP_IP', NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(4, '192.168.88.62', 2, 'IPHONE_14_IP', NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(5, 'device-id-for-the-buried-devices-provider-message', 7, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(6, 'device-id-for-the-rest-provider-message', 7, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(7, 'device-id-for-the-application-message', 7, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(8,'10012db92b',6,NULL,NULL,NULL,'{"Host":"eWeLink_10012db92b","Id":"10012db92b","Ip":"192.168.88.72","Port":"8081","Text":{"apivers":"1","data1":"{\\\"switch\\\":\\\"off\\\",\\\"startup\\\":\\\"off\\\",\\\"pulse\\\":\\\"off\\\",\\\"sledOnline\\\":\\\"on\\\",\\\"fwVersion\\\":\\\"3.6.0\\\",\\\"pulseWidth\\\":500,\\\"rssi\\\":-24}","id":"10012db92b","seq":"73","txtvers":"1","type":"diy_plug"}}',NULL);
INSERT INTO devices VALUES(9,'10011cec96',6,NULL,NULL,NULL,'{"Host":"eWeLink_10011cec96","Id":"10011cec96","Ip":"192.168.88.60","Port":"8081","Text":{"apivers":"1","data1":"{\\\"switch\\\":\\\"off\\\",\\\"startup\\\":\\\"off\\\",\\\"pulse\\\":\\\"off\\\",\\\"sledOnline\\\":\\\"on\\\",\\\"fwVersion\\\":\\\"3.6.0\\\",\\\"pulseWidth\\\":500,\\\"rssi\\\":-69}","id":"10011cec96","seq":"259","txtvers":"1","type":"diy_plug"}}',NULL);
INSERT INTO devices VALUES(10, 'Mhz19Bot', 5, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(11, 'Mhz19ToGoBot', 5, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(12, '18225', 3, NULL, NULL, NULL, NULL, NULL);
INSERT INTO devices VALUES(13, '6613075', 3, NULL, NULL, NULL, NULL, NULL);

COMMIT;