package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var actions = types.ActionImpls{
	types.ACTION_POST_SONOFF_SWITCH_MESSAGE: PostSonoffSwitchMessage,
	types.ACTION_YEELIGHT_DEVICE_SET_POWER:  YeelightDeviceSetPower,
	types.ACTION_MQTT_SET_STATE:             MqttSetState,
	types.ACTION_TELEGRAM_BOT_MESSAGE:       TelegramBotMessage,
	types.ACTION_RECORD_MESSAGE:             RecordMessage,
	types.ACTION_UPSERT_ZIGBEE_DEVICES:      UpsertZigbeeDevices,
	types.ACTION_UPSERT_SONOFF_DEVICE:       UpsertSonoffDevice,
	types.ACTION_PLAY_ALERT:                 PlayAlert,
	types.ACTION_WATCH_CHANGES:              WatchChanges,
}

func Get(fn types.ActionFn) (action types.ActionImpl) {
	action, exist := actions[fn]
	if !exist {
		panic(fmt.Sprintf("Action function %d not yet implemented", fn))
	}
	return
}
