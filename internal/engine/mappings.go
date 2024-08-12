package engine

var Rules = []Rule{

	// test mapping
	{
		Condition: Condition{
			Fn: Equal,
			Args: NamedArgs{
				"Left":  "$deviceClass",
				"Right": DEVICE_CLASS_PINGER,
			},
		},
		Actions: []Action{
			{
				Fn: TelegramBotMessage,
			},
		},
	},

	// balcony ceiling light on/off
	{
		Condition: Condition{
			List: []Condition{
				{
					Fn: Equal,
					Args: NamedArgs{
						"Left":  "$deviceClass",
						"Right": DEVICE_CLASS_ZIGBEE_DEVICE,
					},
				},
				{
					Fn: InList,
					Args: NamedArgs{
						"Value": "$deviceId",
						"List":  []string{"0x00158d0004244bda"},
					},
				},
				{
					Fn: InList,
					Args: NamedArgs{
						"Value": "$message.action",
						"List":  []string{"single_left", "right"},
					},
				},
			},
		},
		Actions: []Action{
			{
				Fn: PostSonoffSwitchMessage,
				Args: NamedArgs{
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
}
