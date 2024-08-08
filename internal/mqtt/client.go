package mqtt

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/registry"
)

var withTag = logger.MakeTag("MQTT")

type service struct {
	out    engine.MessageChan
	client MqttLib.Client
}

var Service engine.Service = &service{}

func (s *service) Receive() engine.MessageChan {
	return s.out
}

func (s *service) Type() engine.ChannelType {
	return engine.CHANNEL_MQTT
}

func (p *parserBase) parse_base() (engine.Message, bool) {

	payload := p.m.Payload()
	topic := p.m.Topic()
	meta := engine.ChannelMeta{MqttTopic: topic}

	outMsg := engine.Message{
		ChannelType: engine.CHANNEL_MQTT,
		ChannelMeta: meta,
		DeviceClass: p.dc,
		Timestamp:   time.Now(),
	}

	tt := strings.Split(strings.TrimLeft(topic, "/"), "/")

	if deviceId := tt[1]; len(tt) >= 2 && p.dc != engine.DEVICE_CLASS_ZIGBEE_BRIDGE {
		outMsg.DeviceId = engine.DeviceId(deviceId)
	}

	if err := json.Unmarshal(payload, &outMsg.Payload); err != nil {
		slog.Warn(withTag("Failed to parse payload as json"), "payload", string(payload[:]), "err", err)
		outMsg.RawPayload = payload
	}

	return outMsg, true
}

func (s *service) Init() {

	s.out = make(engine.MessageChan, 100)

	var handlers = TopicHandlers{
		"zigbee2mqtt/+": func(client MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeDevice(msg).Parse()
			if ok {
				s.out <- outMsg
			}
		},
		"device-pinger/+/status": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewDevicePinger(msg).Parse()
			if ok {
				s.out <- outMsg
			}
		},
		"/VALVE/#": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewValveManipulator(msg).Parse()
			if ok {
				s.out <- outMsg
			}
		},
		"zigbee2mqtt/bridge/devices": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeBridge(msg).Parse()
			if ok {
				s.out <- outMsg
			}
		},
	}

	var defaultMessageHandler = func(client MqttLib.Client, msg MqttLib.Message) {
		slog.Error(withTag("defaultMessageHandler is not expected to be reached"), "topic", msg.Topic())
	}

	var connectHandler = func(client MqttLib.Client) {
		slog.Info(withTag("Connected"), "broker", registry.GetMqttBroker())
	}

	var reconnectHandler = func(client MqttLib.Client, opts *MqttLib.ClientOptions) {
		slog.Info(withTag("Reconnecting..."), "broker", registry.GetMqttBroker())
	}

	var connectLostHandler = func(client MqttLib.Client, err error) {
		slog.Error(withTag("Connection lost"), "error", err)
	}

	opts := MqttLib.NewClientOptions()
	opts.AddBroker(registry.GetMqttBroker())
	opts.SetClientID(registry.Config.MqttClientId)
	opts.SetUsername(registry.Config.MqttUsername)
	opts.SetPassword(registry.Config.MqttPassword)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetAutoReconnect(true)
	opts.OnConnect = connectHandler
	opts.OnReconnecting = reconnectHandler
	opts.OnConnectionLost = connectLostHandler
	s.client = MqttLib.NewClient(opts)
	slog.Debug(withTag("Connecting..."))
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(withTag(""), "error", token.Error())
		return
	}

	for t, h := range handlers {
		s.client.AddRoute(t, h)
		subscribe(s.client, t)
	}
	slog.Debug(withTag("All subscribtions are settled"))

}

func (s *service) Stop() {
	slog.Debug(withTag("Disconnecting..."))
	s.client.Disconnect(250)
}

func subscribe(client MqttLib.Client, topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error(withTag("client.Subscribe()"), "error", token.Error())
	}
	slog.Info(withTag("Subscribed to"), "topic", topic)
}

// var wg sync.WaitGroup
// s.client.AddRoute("zigbee2mqtt/+", zigbeeDeviceHandler)
// s.client.AddRoute("device-pinger/+/status", devicePingerHandler)
// s.client.AddRoute("/VALVE/#", valveManipulatorHandler)
// subscribe_all(s.client, registry.Config.MqttTopics)
// var zigbeeDeviceHandler =
// var devicePingerHandler = func(client MqttLib.Client, msg MqttLib.Message) {
// 	slog.Debug("zigbeeDeviceHandler", "topic", msg.Topic())
// }
// var valveManipulatorHandler = func(client MqttLib.Client, msg MqttLib.Message) {
// 	slog.Debug("zigbeeDeviceHandler", "topic", msg.Topic())
// }
// a, err := AdapterFactory(
// 	msg.Topic(),
// 	msg.Payload(),
// )
// if err != nil {
// 	slog.Error(withTag(err.Error()))
// 	return
// }
// m, err := a.Message()
// if err != nil {
// 	slog.Error(withTag(err.Error()))
// 	return
// }
// s.ch <- m
// func subscribe_all(client MqttLib.Client, topics []string) {
// 	var wg sync.WaitGroup
// 	for _, topic := range topics {
// 		wg.Add(1)
// 		go subscribe(client, topic, &wg)
// 	}
// 	wg.Wait()
// 	slog.Debug(withTag("All subscribtions are settled"))
// }
