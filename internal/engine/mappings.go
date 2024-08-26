package engine

var Rules = []Rule{

	// Comments: "test mapping 1",
	{
		Id:       1,
		Disabled: true,
		Comments: "test mapping 1",
		Condition: Condition{
			Fn: COND_EQUAL,
			Args: Args{
				"Left":  "$deviceClass",
				"Right": DEVICE_CLASS_PINGER,
			},
		},
		Actions: []Action{{
			Fn: ACTION_TELEGRAM_BOT_MESSAGE,
			// Args: Args{"Text": "There is some message from pinger out there"}
		}},
	},

	// Comments: "test mapping for composite condition function",
	{
		Id:       2,
		Disabled: true,
		Comments: "test mapping for composite condition function",
		Condition: Condition{
			List: []Condition{
				{
					Fn: COND_ZIGBEE_DEVICE,
					Args: Args{
						"List": []any{DeviceId("0x00158d0004244bda")},
					},
				},
				{
					Fn: COND_IN_LIST,
					Args: Args{
						"Value": "$message.action",
						"List":  []any{"single_left", "single_right"},
					},
				},
			},
		},
		Actions: []Action{{Fn: ACTION_TELEGRAM_BOT_MESSAGE}},
	},

	// Comments: "balcony ceiling light on/off",
	// 23:44:12.197 DBG [ENGN] New message ChannelType="mqtt (id=1)" ChannelMeta={MqttTopic:zigbee2mqtt/0x00158d0004244bda} DeviceClass="zigbee-device (id=1)" DeviceId=0x00158d0004244bda Payload="map[action:single_right battery:100 device_temperature:30 linkquality:69 power_outage_count:24 voltage:3025]"
	{
		Id:       3,
		Disabled: true,
		Comments: "balcony ceiling light on/off",
		Condition: Condition{
			List: []Condition{
				{
					Fn: COND_EQUAL,
					Args: Args{
						"Left":  "$deviceClass",
						"Right": DEVICE_CLASS_ZIGBEE_DEVICE,
					},
				},
				{
					Fn: COND_IN_LIST,
					Args: Args{
						"Value": "$deviceId",
						"List":  []any{DeviceId("0x00158d0004244bda")},
					},
				},
				{
					Fn: COND_IN_LIST,
					Args: Args{
						"Value": "$message.action",
						"List":  []any{"single_left", "single_right"},
					},
				},
			},
		},
		Actions: []Action{
			{
				Fn: ACTION_POST_SONOFF_SWITCH_MESSAGE,
				Args: Args{
					"Value":    "$message.action",
					"DeviceId": "10011cec96",
				},
				Mapping: Mapping{
					"Value": {
						"single_left":  "on",
						"single_right": "off",
					},
				},
			},
		},
	},

	// Comments: "echo bot",
	{
		Id:       4,
		Disabled: false,
		Comments: "echo bot",
		Condition: Condition{
			Fn: COND_EQUAL,
			Args: Args{
				"Left":  "$deviceClass",
				"Right": DEVICE_CLASS_BOT,
			},
		},
		Actions: []Action{
			{
				Fn:   ACTION_TELEGRAM_BOT_MESSAGE,
				Args: Args{"Text": "Hello!"},
			},
		},
	},
}

// wtf, can we do it more elegant in go? like int(withFn) + int(withList) == 1
// if withList && withFn || (!withList && !withFn) {
// 	panic("unexpected conditions")
// }
// OutChannel: CHANNEL_TELEGRAM,
// nested []model.DbRuleCondition,
// args []model.DbRuleConditionOrActionArgument,
// var nested []model.DbRuleCondition
// var current model.DbRuleCondition
// slices.IndexFunc(conditions);
// if len(ruleall) == 0 {
// 	return
// }
// find root (with null ParentConditionId)
// root, found := lo.Find(conditions, func(c model.DbRuleCondition) bool {
// 	return !c.ParentConditionId.Valid
// })
// children := lo.Filter(conditions, func(c model.DbRuleCondition, i int) bool {
// 	return c.ParentConditionId.Valid && c.ParentConditionId.Int32 == root.Id
// })
// for _, child := range children {
// 	list = append(list, BuildConditions())
// }
// panic("FunctionType for root should not be null")
// nested := lo.Filter(conditions, func(c model.DbRuleCondition, i int) bool {
// 	return c.ParentConditionId.Valid
// })
// for _, c := range conditions {
// 	if c.RuleId == ruleId {
// 	}
// }
// return c_level()
// List := make([]Condition, 0)
// // if con
// for _, c := range conditions {
// 	if c.RuleId != ruleId {
// 		continue
// 	}
// 	if c.FunctionType.Valid {
// 		return Condition{
// 			Fn: CondFn(c.FunctionType.Int32),
// 		}
// 	} else {
// 		// TBD
// 		// list = append(list, Condition{})
// 	}
// }
// return Condition{
// 	List: List,
// }
