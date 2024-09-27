package types

import (
	"encoding/json"
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/logger"
)

type ActionFn byte

var _ fmt.Stringer = (*ActionFn)(nil)
var _ json.Marshaler = (*ActionFn)(nil)

const (
	ACTION_POST_SONOFF_SWITCH_MESSAGE ActionFn = 1
	ACTION_TELEGRAM_BOT_MESSAGE       ActionFn = 2
	ACTION_VALVE_SET_STATE            ActionFn = 3
	ACTION_YEELIGHT_DEVICE_SET_POWER  ActionFn = 4
	ACTION_ZIGBEE2_MQTT_SET_STATE     ActionFn = 5
	ACTION_RECORD_MESSAGE             ActionFn = 6
	ACTION_UPSERT_ZIGBEE_DEVICES      ActionFn = 7
	ACTION_UPSERT_SONOFF_DEVICE       ActionFn = 8
	ACTION_PLAY_ALERT                 ActionFn = 9
)

var ACTION_NAMES = map[ActionFn]string{
	ACTION_POST_SONOFF_SWITCH_MESSAGE: "PostSonoffSwitchMessage",
	ACTION_TELEGRAM_BOT_MESSAGE:       "TelegramBotMessage",
	ACTION_VALVE_SET_STATE:            "ValveSetState",
	ACTION_YEELIGHT_DEVICE_SET_POWER:  "YeelightDeviceSetPower",
	ACTION_ZIGBEE2_MQTT_SET_STATE:     "Zigbee2MqttSetState",
	ACTION_RECORD_MESSAGE:             "RecordMessage",
	ACTION_UPSERT_ZIGBEE_DEVICES:      "UpsertZigbeeDevices",
	ACTION_UPSERT_SONOFF_DEVICE:       "UpsertSonoffDevice",
	ACTION_PLAY_ALERT:                 "PlayAlert",
}

func (fn ActionFn) String() string {
	return ACTION_NAMES[fn]
}

func (fn ActionFn) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, fn)), nil
}

func (fn *ActionFn) UnmarshalJSON(b []byte) (err error) {
	var v any
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for action, name := range ACTION_NAMES {
		switch vtyped := v.(type) {
		case string:
			if name == vtyped {
				*fn = action
				return
			}
		case float64:
			if float64(action) == vtyped {
				*fn = ActionFn(vtyped)
				return
			}
		}
	}
	return fmt.Errorf("failed to unmarshal %v (type=%T) to ActionFn", v, v)
}

type ActionImpl func(messages []Message, args Args, mapping Mapping, engine EngineAsSupplier, tag logger.Tag) error

type ActionImpls map[ActionFn]ActionImpl
