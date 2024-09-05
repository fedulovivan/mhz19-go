package engine

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

func GetStaticRules() []types.Rule {
	return []types.Rule{

		// system rule to save received message in db
		{
			Id:       1,
			Disabled: false,
			Name:     "system rule to save received message in db",
			Condition: types.Condition{
				List: []types.Condition{
					{
						Fn: types.COND_NOT_EQUAL,
						Args: types.Args{
							"Left":  "$deviceClass",
							"Right": types.DEVICE_CLASS_ZIGBEE_BRIDGE,
						},
					},
					{
						Fn: types.COND_NOT_EQUAL,
						Args: types.Args{
							"Left":  "$channelType",
							"Right": types.CHANNEL_SONOFF,
						},
					},
				},
			},
			Actions: []types.Action{{Fn: types.ACTION_RECORD_MESSAGE}},
		},

		// system rule to create devices upon receiving message from bridge
		{
			Id:       2,
			Disabled: false,
			Name:     "system rule to create devices upon receiving message from bridge",
			Condition: types.Condition{
				Fn: types.COND_EQUAL,
				Args: types.Args{
					"Left":  "$deviceClass",
					"Right": types.DEVICE_CLASS_ZIGBEE_BRIDGE,
				},
			},
			Actions: []types.Action{{Fn: types.ACTION_UPSERT_ZIGBEE_DEVICES}},
		},

		// Comments: "test mapping 1",
		{
			Id:       3,
			Disabled: true,
			Name:     "test mapping 1",
			Condition: types.Condition{
				Fn: types.COND_EQUAL,
				Args: types.Args{
					"Left":  "$deviceClass",
					"Right": types.DEVICE_CLASS_PINGER,
				},
			},
			Actions: []types.Action{{
				Fn: types.ACTION_TELEGRAM_BOT_MESSAGE,
				// Args: types.Args{"Text": "There is some message from pinger out there"}
			}},
		},

		// Comments: "test mapping for composite condition function",
		{
			Id:       4,
			Disabled: true,
			Name:     "test mapping for composite condition function",
			Condition: types.Condition{
				List: []types.Condition{
					{
						Fn: types.COND_ZIGBEE_DEVICE,
						Args: types.Args{
							"List": []any{types.DeviceId("0x00158d0004244bda")},
						},
					},
					{
						Fn: types.COND_IN_LIST,
						Args: types.Args{
							"Value": "$message.action",
							"List":  []any{"single_left", "single_right"},
						},
					},
				},
			},
			Actions: []types.Action{{Fn: types.ACTION_TELEGRAM_BOT_MESSAGE}},
		},

		// Comments: "balcony ceiling light on/off",
		// 23:44:12.197 DBG [ENGN] New message ChannelType="mqtt (id=1)" ChannelMeta={MqttTopic:zigbee2mqtt/0x00158d0004244bda} DeviceClass="zigbee-device (id=1)" DeviceId=0x00158d0004244bda Payload="map[action:single_right battery:100 device_temperature:30 linkquality:69 power_outage_count:24 voltage:3025]"
		{
			Id:       5,
			Disabled: true,
			Name:     "balcony ceiling light on/off",
			Condition: types.Condition{
				List: []types.Condition{
					{
						Fn: types.COND_ZIGBEE_DEVICE,
						Args: types.Args{
							"List": []any{types.DeviceId("0x00158d0004244bda")},
						},
					},
					{
						Fn: types.COND_IN_LIST,
						Args: types.Args{
							"Value": "$message.action",
							"List":  []any{"single_left", "single_right"},
						},
					},
				},
			},
			Actions: []types.Action{
				{
					Fn: types.ACTION_POST_SONOFF_SWITCH_MESSAGE,
					Args: types.Args{
						"Command":  "$message.action",
						"DeviceId": types.DeviceId("10011cec96"),
					},
					Mapping: types.Mapping{
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
			Id:       6,
			Disabled: true,
			Name:     "echo bot",
			// Throttle: time.Second / 2,
			Condition: types.Condition{
				Fn: types.COND_EQUAL,
				Args: types.Args{
					"Left":  "$deviceClass",
					"Right": types.DEVICE_CLASS_BOT,
				},
			},
			Actions: []types.Action{
				{
					Fn: types.ACTION_TELEGRAM_BOT_MESSAGE,
					Args: types.Args{
						"Text": "accumulated {{len .Messages}} message(s):\n{{ range .Messages }}message from {{ .DeviceId }} with text {{ .Payload.Text }}\n{{ end }}",
					},
				},
			},
		},
	}
}
