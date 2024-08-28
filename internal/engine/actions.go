package engine

import (
	"encoding/json"
	"log/slog"

	"fmt"
)

type ActionFn byte

type GetProvider func(ch ChannelType) Provider

type ActionImpl func(mm []Message, a Action, gs GetProvider, e *engine)

type ActionImpls map[ActionFn]ActionImpl

const (
	ACTION_POST_SONOFF_SWITCH_MESSAGE ActionFn = 1
	ACTION_TELEGRAM_BOT_MESSAGE       ActionFn = 2
	ACTION_VALVE_SET_STATE            ActionFn = 3
	ACTION_YEELIGHT_DEVICE_SET_POWER  ActionFn = 4
	ACTION_ZIGBEE2_MQTT_SET_STATE     ActionFn = 5
	ACTION_RECORD_MESSAGE             ActionFn = 6
)

var ACTION_NAMES = map[ActionFn]string{
	ACTION_POST_SONOFF_SWITCH_MESSAGE: "PostSonoffSwitchMessage",
	ACTION_TELEGRAM_BOT_MESSAGE:       "TelegramBotMessage",
	ACTION_VALVE_SET_STATE:            "ValveSetState",
	ACTION_YEELIGHT_DEVICE_SET_POWER:  "YeelightDeviceSetPower",
	ACTION_ZIGBEE2_MQTT_SET_STATE:     "Zigbee2MqttSetState",
	ACTION_RECORD_MESSAGE:             "RecordMessage",
}

func (s ActionFn) String() string {
	return fmt.Sprintf("%v (id=%d)", ACTION_NAMES[s], s)
}

// func (s *ActionFn) MarshalJSON() ([]byte, error) {
// 	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
// }

var PostSonoffSwitchMessage ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	panic("not yet implemented")
}

var YeelightDeviceSetPower ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	panic("not yet implemented")
}

var Zigbee2MqttSetState ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	// s := gs(CHANNEL_MQTT)
	// s.Send("foo1")
	panic("not yet implemented")
}

var ValveSetState ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	panic("not yet implemented")
}

var TelegramBotMessage ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	p := gs(CHANNEL_TELEGRAM)
	text := a.Args["Text"]
	if text != nil {
		p.Send(text)
	} else {
		p.Send(json.Marshal(mm[0]))
	}
}

var RecordMessage ActionImpl = func(mm []Message, a Action, gs GetProvider, e *engine) {
	err := e.options.messageService.Create(mm[0])
	if err != nil {
		slog.Error(err.Error())
	}
}

var actions = ActionImpls{
	ACTION_POST_SONOFF_SWITCH_MESSAGE: PostSonoffSwitchMessage,
	ACTION_YEELIGHT_DEVICE_SET_POWER:  YeelightDeviceSetPower,
	ACTION_ZIGBEE2_MQTT_SET_STATE:     Zigbee2MqttSetState,
	ACTION_VALVE_SET_STATE:            ValveSetState,
	ACTION_TELEGRAM_BOT_MESSAGE:       TelegramBotMessage,
	ACTION_RECORD_MESSAGE:             RecordMessage,
}
