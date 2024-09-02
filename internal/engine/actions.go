package engine

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/devices"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type ActionFn byte

type GetProviderFn func(ch types.ChannelType) ChannelProvider

type ActionImpl func(mm []types.Message, a Action, e *engine)

type ActionImpls map[ActionFn]ActionImpl

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

var PostSonoffSwitchMessage ActionImpl = func(mm []types.Message, a Action, e *engine) {

	areader := NewArgReader(mm[0], a.Args, a.Mapping)
	cmd := areader.Get("Command")
	deviceId := areader.Get("DeviceId")
	if !areader.Ok() {
		slog.Error(areader.Error().Error())
		return
	}

	device, err := e.options.devicesService.GetOne(deviceId.(types.DeviceId))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	djson := gabs.Wrap(device.Json)
	ip := djson.Path("ip").Data().(string)
	port := djson.Path("port").Data().(string)

	url := fmt.Sprintf("http://%v:%v/zeroconf/switch", ip, port)
	payload := []byte(fmt.Sprintf(`{"data":{"switch":"%v"}}`, cmd))
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err == nil && res.StatusCode == 200 {
		fmt.Println("success")
	}
	fmt.Println(url, string(payload), res, err)
}

var YeelightDeviceSetPower ActionImpl = func(mm []types.Message, a Action, e *engine) {
	panic("not yet implemented")
}

var Zigbee2MqttSetState ActionImpl = func(mm []types.Message, a Action, e *engine) {
	panic("not yet implemented")
}

var ValveSetState ActionImpl = func(mm []types.Message, a Action, e *engine) {
	panic("not yet implemented")
}

var TelegramBotMessage ActionImpl = func(mm []types.Message, a Action, e *engine) {
	p := e.getPrivider(types.CHANNEL_TELEGRAM)
	text := a.Args["Text"]
	if text != nil {
		p.Send(text)
	} else {
		p.Send(json.Marshal(mm[0]))
	}
}

var RecordMessage ActionImpl = func(mm []types.Message, a Action, e *engine) {
	err := e.options.messageService.Create(mm[0])
	if err != nil {
		slog.Error(e.options.logTag(err.Error()))
	}
}

// system action to create devices upon receiving message from zigbee2mqtt bridge
// see https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html#zigbee2mqtt-bridge-devices
// and json example at assets/bridge-devices-message.json
var UpsertZigbeeDevices ActionImpl = func(mm []types.Message, a Action, e *engine) {
	devicesjson := gabs.Wrap(mm[0].Payload)
	out := make([]devices.Device, 0)
	for _, d := range devicesjson.Children() {
		dtype := d.Path("type").Data().(string)
		if dtype == "Coordinator" {
			continue
		}
		out = append(out, devices.Device{
			DeviceClassId: types.DEVICE_CLASS_ZIGBEE_DEVICE,
			DeviceId:      types.DeviceId(d.Path("ieee_address").Data().(string)),
			Comments:      d.Path("definition.description").Data().(string),
			Origin:        "bridge-upsert",
			Json:          d.Data(),
		})
	}
	err := e.options.devicesService.Upsert(out)
	if err != nil {
		slog.Error(e.options.logTag(err.Error()))
	}
}

var actions = ActionImpls{
	ACTION_POST_SONOFF_SWITCH_MESSAGE: PostSonoffSwitchMessage,
	ACTION_YEELIGHT_DEVICE_SET_POWER:  YeelightDeviceSetPower,
	ACTION_ZIGBEE2_MQTT_SET_STATE:     Zigbee2MqttSetState,
	ACTION_VALVE_SET_STATE:            ValveSetState,
	ACTION_TELEGRAM_BOT_MESSAGE:       TelegramBotMessage,
	ACTION_RECORD_MESSAGE:             RecordMessage,
	ACTION_UPSERT_ZIGBEE_DEVICES:      UpsertZigbeeDevices,
}
