- No new messages after mqtt disconnect/autoreconnect (`Connection lost error="pingresp not received, disconnecting"` and later `Connected broker=tcp://macmini:1883`) + no reconnection if network was not available on app startup and returned online later
- Mappings could be built-in (system) and loaded from db (user)
- create test service for sonoff wifi devices (poll them periodically to receive status updates)
- refactor to avoid "startup messages" in unit tests
- add support of hi-level composite functions like ZigbeeDevice under the hood consisting of Equal and InList
- Service.Receive() schould be placed in "base" struct
- find out why "make test" and "vscode" report different coverage statistics: 86.9% vs 100%
- think how to revive Service.Type() which is looking redundant now
- (+) wrap internal/mqtt/client.go and internal/tbot/tbot.go into structs
- (+) for mqtt client, rather than hardocding in defaultMessageHandler, define rules/adapters for transforming topic and payload into final message per device class
- (+) consider replacing adapters with mqttClient.AddRoute(), also add warning for messages captured by defaultMessageHandler (assuming all topics we subscribe should have own handlers and default one should not be reached)


### new mapping rule structure
```golang
    type ConditionArgument interface {
        int | string | bool
    }
    type LogicOp string
    type ConditionFn string
    type ActionFn string
    const AND ConditionOp = "AND"
    const OR ConditionOp = "OR"
    type Action struct {
        Fn ActionFn
        Args []ConditionArgument
        Mapping
    }
    type Condition struct {
        Fn ConditionFn
        Args []ConditionArgument
        LogicOp LogicOp
        List []Condition
    }
    type Rule struct {
        Condition Condition
        Actions []Action
        Throttle int
    }    
```

// example 1
// execute RelayToTelegram action for all messages with deviceId=192.168.88.1
// (note that channel and device class are ignored)
{
    Condition{ Fn: "Equals" Args: ["$deviceId", "192.168.88.1"] }
    Actions{ RelayToTelegram }
}

// example 2
// execute RelayToTelegram, SaveToSqlite actions for all messages
// which conform
// deviceClass=ZIGBEE AND deviceId=0x00158d0000c2fa6e AND (Changed message.state OR NotNil message.lastSeen)
{
    Condition{Op AND, List{
        { Fn: "Equals" Args: ["$deviceClass", ZIGBEE] }
        { Fn: "Equals" Args: ["$deviceId", "0x00158d0000c2fa6e"] }
        { Op: OR List{
            { Fn: "Changed" Args: ["$message.state"] }
            { Fn: "NotNil" Args: ["$message.lastSeen"] }
        } }
    }}
    Actions{ RelayToTelegram, SaveToSqlite }
}

// example 3
// implement high-level helper functions
// func ZigbeeDevice () Equals deviceClass AND  
{
    Condition{ Fn: "ZigbeeDevice" Args: ["0x00158d0000c2fa6e", "0x00158d000405811b"] }   
}

// rules
//   id int
//   comments string
//   enabled bool
//   throttle int
// rule_conditions
//   id
//   rule_id
//   function_name
//   logical_operation
//   parent_condition_id
// rule_actions
//   id
//   rule_id
//   function_name
//   device_id
// rule_condition_and_action_arguments
//   id
//   rule_id
//   condition_id
//   action_id
//   string_value
//   device_id_value
// rule_action_mappings
//   id
//   rule_id
//   action_id
//   key
//   value