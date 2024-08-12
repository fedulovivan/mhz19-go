CREATE TABLE devices (
   	id INTEGER PRIMARY KEY AUTOINCREMENT,
    orig_id TEXT NOT NULL,
    device_class TEXT NOY NULL,
	name TEXT,
	comments TEXT,
);

CREATE TABLE rules (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	comments TEXT NOT NULL,
	is_enabled INTEGER DEFAULT 1 NOT NULL,
	throttle INTEGER
);

CREATE TABLE rule_conditions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_name TEXT,
	logical_operation TEXT,
	parent_condition_id INTEGER,
	CONSTRAINT rule_conditions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT rule_conditions_fk_parent FOREIGN KEY (parent_condition_id) REFERENCES rule_conditions(id)
);

CREATE TABLE rule_actions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	rule_id INTEGER NOT NULL,
	function_name TEXT NOT NULL,
    device_id INTEGER,
    CONSTRAINT rule_actions_fk_rules FOREIGN KEY (rule_id) REFERENCES rules(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_actions_fk_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_condition_or_action_arguments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    argument_name TEXT NOT NULL,
    condition_id INTEGER,
    action_id INTEGER,
    string_value TEXT,
    device_id INTEGER,
    CONSTRAINT rule_ca_arguments_fk_conditions FOREIGN KEY (condition_id) REFERENCES rule_conditions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_actions FOREIGN KEY (action_id) REFERENCES rule_actions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT rule_ca_arguments_fk_devices FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE rule_action_argument_mappings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    argument_id INTEGER,
    key TEXT,
    value TEXT,
    CONSTRAINT rule_action_mappings_fk_arguments FOREIGN KEY (argument_id) REFERENCES rule_condition_or_action_arguments(id) ON DELETE CASCADE ON UPDATE CASCADE
);