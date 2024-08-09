package engine

import (
	"fmt"
	"time"
)

type JsonPayload interface{}

// type JsonPayload map[string]interface{}

type DeviceId string

type DeviceClass byte

const (
	DEVICE_CLASS_ZIGBEE_DEVICE DeviceClass = 1
	DEVICE_CLASS_PINGER        DeviceClass = 2
	DEVICE_CLASS_VALVE         DeviceClass = 3
	DEVICE_CLASS_ZIGBEE_BRIDGE DeviceClass = 4
	DEVICE_CLASS_BOT           DeviceClass = 5
)

var DEVICE_CLASS_NAMES = map[DeviceClass]string{
	DEVICE_CLASS_ZIGBEE_DEVICE: "zigbee-device",
	DEVICE_CLASS_PINGER:        "device-pinger",
	DEVICE_CLASS_VALVE:         "valve-manipulator",
	DEVICE_CLASS_ZIGBEE_BRIDGE: "zigbee-bridge",
	DEVICE_CLASS_BOT:           "telegram-bot",
}

func (s DeviceClass) String() string {
	if s == 0 {
		return "<unknown>"
	}
	return fmt.Sprintf("%v (id=%d)", DEVICE_CLASS_NAMES[s], s)
}

type ChannelType byte

const (
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
}

func (s ChannelType) String() string {
	if s == 0 {
		return "<unknown>"
	}
	return fmt.Sprintf("%v (id=%d)", CHANNEL_NAMES[s], s)
}

type ChannelMeta struct {
	MqttTopic string
}

type Message struct {
	// additional metadata specific for the current channel
	ChannelMeta ChannelMeta
	// channel, which was used to receive message
	ChannelType ChannelType `json:"channelType"`
	// device class, see DeviceClass
	DeviceClass DeviceClass `json:"deviceClass"`
	// device id extracted from topic
	DeviceId DeviceId `json:"deviceId"`
	// time when message was received by backend
	Timestamp time.Time `json:"timestamp"`
	// parsed message payload json
	Payload JsonPayload `json:"payload"`
	// filled only if failed to parse into json
	RawPayload []byte `json:"rawPayload"`
}

type MessageChan chan Message

type Service interface {
	Receive() MessageChan
	Type() ChannelType
	Init()
	Stop()
}

type Arg interface{}
type Mapping map[string]string

type Op string
type CondFn string
type ActionFn string

const (
	And Op = "And"
	Or  Op = "Or"
)

const (
	ZigbeeDevice CondFn = "ZigbeeDevice"
	Equal        CondFn = "Equal"
	NotEqual     CondFn = "NotEqual"
	InList       CondFn = "InList"
	Changed      CondFn = "Changed"
	NotNil       CondFn = "NotNil"
)

const (
	PostSonoffSwitchMessage ActionFn = "PostSonoffSwitchMessage"
	YeelightDeviceSetPower  ActionFn = "YeelightDeviceSetPower"
	Zigbee2MqttSetState     ActionFn = "Zigbee2MqttSetState"
	ValveSetState           ActionFn = "ValveSetState"
	TelegramBotMessage      ActionFn = "TelegramBotMessage"
)

type Action struct {
	DeviceId string
	Fn       ActionFn
	Args     []Arg
	Mapping
}

type Condition struct {
	Fn   CondFn
	Args []Arg
	Op   Op
	List []Condition
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

type Rule struct {
	Throttle  time.Duration
	Condition Condition
	Actions   []Action
}

var Rules = []Rule{
	// balcony ceiling light
	{
		Condition: Condition{
			Op: And,
			List: []Condition{
				{
					Fn: ZigbeeDevice,
					Args: []Arg{
						"0x00158d0004244bda",
					},
				},
				{
					Fn: InList,
					Args: []Arg{
						"$message.action",
						"single_left",
						"right",
					},
				},
			},
		},
		Actions: []Action{
			{
				Fn:       PostSonoffSwitchMessage,
				DeviceId: "10011cec96",
				Args: []Arg{
					"$message.action",
				},
				Mapping: Mapping{
					"single_left":  "on",
					"single_right": "off",
				},
			},
		},
	},
}

func init() {
	fmt.Println(Rules)
}
