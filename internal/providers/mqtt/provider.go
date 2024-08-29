package mqtt_service

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
)

var logTag = logger.MakeTag(logger.MQTT)

// implements [engine.Provider]
type provider struct {
	engine.ProviderBase
	client MqttLib.Client
}

var Provider engine.Provider = &provider{}

func (s *provider) Channel() engine.ChannelType {
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
		slog.Warn(logTag("Failed to parse payload as json"), "payload", string(payload[:]), "err", err)
		outMsg.RawPayload = payload
	}

	return outMsg, true
}

func (s *provider) Init() {

	s.Out = make(engine.MessageChan, 100)

	var handlers = TopicHandlers{
		"zigbee2mqtt/+": func(client MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeDevice(msg).Parse()
			if ok {
				s.Out <- outMsg
			}
		},
		"device-pinger/+/status": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewDevicePinger(msg).Parse()
			if ok {
				s.Out <- outMsg
			}
		},
		"/VALVE/#": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewValveManipulator(msg).Parse()
			if ok {
				s.Out <- outMsg
			}
		},
		"zigbee2mqtt/bridge/devices": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeBridge(msg).Parse()
			if ok {
				s.Out <- outMsg
			}
		},
	}

	var defaultMessageHandler = func(client MqttLib.Client, msg MqttLib.Message) {
		slog.Error(logTag("defaultMessageHandler is not expected to be reached"), "topic", msg.Topic())
	}

	var connectHandler = func(client MqttLib.Client) {
		slog.Info(logTag("Connected"), "broker", app.GetMqttBroker())
		for t := range handlers {
			subscribe(client, t)
		}
		slog.Debug(logTag("All subscribtions are settled"))
	}

	var reconnectHandler = func(client MqttLib.Client, opts *MqttLib.ClientOptions) {
		slog.Warn(logTag("Reconnecting..."), "broker", app.GetMqttBroker())
	}

	var connectLostHandler = func(client MqttLib.Client, err error) {
		slog.Error(logTag("Connection lost"), "error", err)
	}

	// build opts
	opts := MqttLib.NewClientOptions()
	opts.AddBroker(app.GetMqttBroker())
	opts.SetClientID(app.Config.MqttClientId)
	opts.SetUsername(app.Config.MqttUsername)
	opts.SetPassword(app.Config.MqttPassword)
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetReconnectingHandler(reconnectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)

	// attach logger
	if app.Config.MqttDebug {
		sloghandler := slog.Default().Handler()
		// fmt.Printf("%T", sloghandler)
		MqttLib.ERROR = slog.NewLogLogger(sloghandler, slog.LevelError)
		MqttLib.CRITICAL = slog.NewLogLogger(sloghandler, slog.LevelError)
		MqttLib.WARN = slog.NewLogLogger(sloghandler, slog.LevelWarn)
		MqttLib.DEBUG = slog.NewLogLogger(sloghandler, slog.LevelDebug)
	}

	// create client
	s.client = MqttLib.NewClient(opts)

	// register routes
	for t, h := range handlers {
		s.client.AddRoute(t, h)
	}

	// connect
	slog.Debug(logTag("Connecting..."))
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(logTag("Initial connect"), "error", token.Error())
	}

}

func (s *provider) Stop() {
	slog.Debug(logTag("Disconnecting..."))
	if s.client.IsConnected() {
		s.client.Disconnect(250)
	} else {
		slog.Warn(logTag("Not connected"))
	}
}

func subscribe(client MqttLib.Client, topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error(logTag("client.Subscribe()"), "error", token.Error())
	}
	slog.Info(logTag("Subscribed to"), "topic", topic)
}
