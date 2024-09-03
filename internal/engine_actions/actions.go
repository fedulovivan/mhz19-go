package engine_actions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Actions = types.ActionImpls{
	types.ACTION_POST_SONOFF_SWITCH_MESSAGE: PostSonoffSwitchMessage,
	types.ACTION_YEELIGHT_DEVICE_SET_POWER:  YeelightDeviceSetPower,
	types.ACTION_ZIGBEE2_MQTT_SET_STATE:     Zigbee2MqttSetState,
	types.ACTION_VALVE_SET_STATE:            ValveSetState,
	types.ACTION_TELEGRAM_BOT_MESSAGE:       TelegramBotMessage,
	types.ACTION_RECORD_MESSAGE:             RecordMessage,
	types.ACTION_UPSERT_ZIGBEE_DEVICES:      UpsertZigbeeDevices,
}
