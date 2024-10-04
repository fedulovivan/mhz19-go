package mqtt_provider

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var tag = logger.NewTag(logger.MQTT)

type provider struct {
	engine.ProviderBase
	client MqttLib.Client
}

func NewProvider() types.ChannelProvider {
	return &provider{
		ProviderBase: engine.ProviderBase{
			MessagesChan: make(types.MessageChan, 100),
		},
	}
}

func (p *provider) Send(a ...any) (err error) {
	topic, ok1 := a[0].(string)
	payload, ok2 := a[1].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("Send() expects two string arguments: topic and payload")
	}
	token := p.client.Publish(
		topic,
		0,
		false,
		payload,
	)
	token.Wait()
	return token.Error()
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_MQTT
}

func (p *parserBase) parse_base() (types.Message, bool) {

	payload := p.m.Payload()
	topic := p.m.Topic()
	meta := types.ChannelMeta{MqttTopic: topic}

	outMsg := types.Message{
		Id:            types.MessageIdSeq.Inc(),
		Timestamp:     time.Now(),
		ChannelType:   types.CHANNEL_MQTT,
		ChannelMeta:   &meta,
		DeviceClass:   p.dc,
		FromEndDevice: false,
	}

	tt := strings.Split(strings.TrimLeft(topic, "/"), "/")

	if deviceId := tt[1]; len(tt) >= 2 {
		outMsg.DeviceId = types.DeviceId(deviceId)
	}

	if err := json.Unmarshal(payload, &outMsg.Payload); err != nil {
		slog.Warn(tag.F("Failed to parse payload as json"), "payload", string(payload[:]), "err", err)
		outMsg.RawPayload = payload
	}

	return outMsg, true
}

func (p *provider) Init() {
	var handlers = TopicHandlers{
		"zigbee2mqtt/+": func(client MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeDevice(msg).Parse()
			if ok {
				p.Push(outMsg)
			}
		},
		"device-pinger/+/status": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewDevicePinger(msg).Parse()
			if ok {
				p.Push(outMsg)
			}
		},
		"/VALVE/+/STATE/STATUS": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewValveManipulator(msg).Parse()
			if ok {
				p.Push(outMsg)
			}
		},
		"zigbee2mqtt/bridge/devices": func(c MqttLib.Client, msg MqttLib.Message) {
			outMsg, ok := NewZigbeeBridge(msg).Parse()
			if ok {
				p.Push(outMsg)
			}
		},
	}

	var defaultMessageHandler = func(client MqttLib.Client, msg MqttLib.Message) {
		slog.Error(tag.F("defaultMessageHandler is not expected to be reached"), "topic", msg.Topic())
		counters.Inc(counters.ERRORS_ALL)
	}

	var connectHandler = func(client MqttLib.Client) {
		slog.Info(tag.F("Connected"), "broker", app.GetMqttBrokerUrl())
		for t := range handlers {
			subscribe(client, t)
		}
		slog.Debug(tag.F("All subscriptions are settled"))
	}

	var reconnectHandler = func(client MqttLib.Client, opts *MqttLib.ClientOptions) {
		slog.Warn(tag.F("Reconnecting..."), "broker", app.GetMqttBrokerUrl())
	}

	var connectLostHandler = func(client MqttLib.Client, err error) {
		slog.Error(tag.F("Connection lost"), "error", err)
		counters.Inc(counters.ERRORS_ALL)
	}

	// build opts
	opts := MqttLib.NewClientOptions()
	opts.AddBroker(app.GetMqttBrokerUrl())
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
		MqttLib.ERROR = slog.NewLogLogger(sloghandler, slog.LevelError)
		MqttLib.CRITICAL = slog.NewLogLogger(sloghandler, slog.LevelError)
		MqttLib.WARN = slog.NewLogLogger(sloghandler, slog.LevelWarn)
		MqttLib.DEBUG = slog.NewLogLogger(sloghandler, slog.LevelDebug)
	}

	// create client
	p.client = MqttLib.NewClient(opts)

	// register routes
	for t, h := range handlers {
		p.client.AddRoute(t, h)
	}

	// connect
	slog.Debug(tag.F("Connecting..."))
	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error(tag.F("Initial connect"), "error", token.Error())
		counters.Inc(counters.ERRORS_ALL)
	}

}

func (p *provider) Stop() {
	slog.Debug(tag.F("Disconnecting..."))
	if p.client.IsConnected() {
		p.client.Disconnect(250)
	} else {
		slog.Warn(tag.F("Not connected"))
	}
}

func subscribe(client MqttLib.Client, topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error(tag.F("client.Subscribe()"), "error", token.Error())
		counters.Inc(counters.ERRORS_ALL)
	}
	slog.Info(tag.F("Subscribed to"), "topic", topic)
}
