package types

import "fmt"

type ActionFn byte

const (
	ACTION_POST_SONOFF_SWITCH_MESSAGE ActionFn = 1
	ACTION_TELEGRAM_BOT_MESSAGE       ActionFn = 2
	ACTION_VALVE_SET_STATE            ActionFn = 3
	ACTION_YEELIGHT_DEVICE_SET_POWER  ActionFn = 4
	ACTION_ZIGBEE2_MQTT_SET_STATE     ActionFn = 5
	ACTION_RECORD_MESSAGE             ActionFn = 6
	ACTION_UPSERT_ZIGBEE_DEVICES      ActionFn = 7
)

var ACTION_NAMES = map[ActionFn]string{
	ACTION_POST_SONOFF_SWITCH_MESSAGE: "PostSonoffSwitchMessage",
	ACTION_TELEGRAM_BOT_MESSAGE:       "TelegramBotMessage",
	ACTION_VALVE_SET_STATE:            "ValveSetState",
	ACTION_YEELIGHT_DEVICE_SET_POWER:  "YeelightDeviceSetPower",
	ACTION_ZIGBEE2_MQTT_SET_STATE:     "Zigbee2MqttSetState",
	ACTION_RECORD_MESSAGE:             "RecordMessage",
	ACTION_UPSERT_ZIGBEE_DEVICES:      "UpsertZigbeeDevices",
}

func (s ActionFn) String() string {
	return fmt.Sprintf("%v (id=%d)", ACTION_NAMES[s], s)
}

func (s *ActionFn) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, ACTION_NAMES[*s])), nil
}

type ActionImpl func(messages []Message, action Action, engine EngineAsSupplier) error

type ActionImpls map[ActionFn]ActionImpl
