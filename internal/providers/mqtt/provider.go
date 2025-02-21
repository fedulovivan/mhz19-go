package mqtt_provider

import (
	"fmt"
	"log/slog"

	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var tag = utils.NewTag(logger.MQTT)

type TopicHandlers map[string]MqttLib.MessageHandler

type provider struct {
	engine.ProviderBase
	client MqttLib.Client
}

var _ types.ChannelProvider = (*provider)(nil)

func NewProvider() *provider {
	return &provider{}
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

func (p *provider) Init() {

	p.ProviderBase.Init()

	var handlers = TopicHandlers{
		"zigbee2mqtt/+": func(client MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_ZIGBEE_DEVICE, true, 1))
		},
		"device-pinger/+/status": func(c MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_PINGER, true, 1))
		},
		"/VALVE/+/STATE/STATUS": func(c MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_VALVE, true, 1))
		},
		"valves-manipulator/+/status": func(c MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_VALVE, true, 1))
		},
		"zigbee2mqtt/bridge/devices": func(c MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_ZIGBEE_BRIDGE, false, 1))
		},
		"espresense/devices/+/+": func(c MqttLib.Client, msg MqttLib.Message) {
			p.Push(Parse(msg, types.DEVICE_CLASS_ESPRESENCE_DEVICE, true, 2))
		},
	}

	var defaultMessageHandler = func(client MqttLib.Client, msg MqttLib.Message) {
		slog.Error(tag.F("defaultMessageHandler is not expected to be reached"), "topic", msg.Topic())
		counters.Inc(counters.ERRORS_ALL)
		counters.Errors.WithLabelValues(logger.MOD_MQTT).Inc()
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
		counters.Errors.WithLabelValues(logger.MOD_MQTT).Inc()
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
	opts.SetConnectRetry(true)

	// attach logger
	if app.Config.MqttDebug {
		prefix := string(logger.MQTT) + " "
		sloghandler := slog.Default().Handler()
		lerror := slog.NewLogLogger(sloghandler, slog.LevelError)
		lcritical := slog.NewLogLogger(sloghandler, slog.LevelError)
		lwarn := slog.NewLogLogger(sloghandler, slog.LevelWarn)
		ldebug := slog.NewLogLogger(sloghandler, slog.LevelDebug)
		lerror.SetPrefix(prefix)
		lcritical.SetPrefix(prefix)
		lwarn.SetPrefix(prefix)
		ldebug.SetPrefix(prefix)
		MqttLib.ERROR = lerror
		MqttLib.CRITICAL = lcritical
		MqttLib.WARN = lwarn
		MqttLib.DEBUG = ldebug
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
		counters.Errors.WithLabelValues(logger.MOD_MQTT).Inc()
	}

}

func (p *provider) Stop() {
	slog.Debug(tag.F("Disconnecting..."))
	if p.client.IsConnected() {
		p.client.Disconnect(250)
	} else {
		slog.Warn(tag.F("Not connected"))
	}
	p.ProviderBase.Stop()
}

func subscribe(client MqttLib.Client, topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error(tag.F("client.Subscribe()"), "error", token.Error())
		counters.Inc(counters.ERRORS_ALL)
		counters.Errors.WithLabelValues(logger.MOD_MQTT).Inc()
	}
	slog.Info(tag.F("Subscribed to"), "topic", topic)
}

// "zigbee2mqtt/bridge/event": func(c MqttLib.Client, msg MqttLib.Message) {
// 	outMsg, ok := NewZigbeeBridge(msg).Parse()
// 	if ok {
// 		p.Push(outMsg)
// 	}
// },

// ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
// defer cancel()
// p.initialConnectWithRetry(ctx)

// func (p *provider) initialConnectWithRetry(ctx context.Context) {
// 	retries := 0
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			slog.Error(tag.F("Initial connect failed"), "retries", retries)
// 			slog.Error(tag.F("context cancelled"), "err", ctx.Err())
// 			counters.Inc(counters.ERRORS_ALL)
// 			counters.Errors.WithLabelValues(logger.MOD_MQTT).Inc()
// 			return
// 		default:
// 			slog.Debug(tag.F("Trying to connect..."))
// 			if token := p.client.Connect(); token.Wait() && token.Error() != nil {
// 				slog.Error(tag.F("initialConnectWithRetry"), "err", token.Error())
// 				slog.Warn(tag.F("Retrying in a second..."))
// 				time.Sleep(time.Second)
// 				retries += 1
// 			} else {
// 				return
// 			}
// 		}
// 	}
// }
