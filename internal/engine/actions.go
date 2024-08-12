package engine

type ActionFnName string

type Mapping map[string](map[string]string)

const (
	PostSonoffSwitchMessage ActionFnName = "PostSonoffSwitchMessage"
	YeelightDeviceSetPower  ActionFnName = "YeelightDeviceSetPower"
	Zigbee2MqttSetState     ActionFnName = "Zigbee2MqttSetState"
	ValveSetState           ActionFnName = "ValveSetState"
	TelegramBotMessage      ActionFnName = "TelegramBotMessage"
)
