### Project «mhz19-go»

This project is backend for the home automation server, written on golang. Evolution of another myne project [mhz19-next](https://github.com/fedulovivan/mhz19-next) which was a typescript-based. 

### Real use cases

- Switch smart ceiling light on/off upon receiving message from smart wall switch
- Automatically switch storage room light upon receiving message from movement sensor
- Automatically switch storage room ventilation, when movement sensor reports a man presense and room door sensor reports it is closed
- Play loud alert sound, notify owner via telegram and cut off home water supply upon receiving message from one of water leakage sensors
- Notify owner then some guarded door (equipped with smart sensor) was opened/closed and user is not at home

### Project goal

- Create pure offline, vendor agnostic and fully controlled local home automation server
- Bring project to modern stack and learn golang while project migration from typescript

### Architecture

- Channel providers which collect messages from different channels and devices into unified stream. Supported channels are mqtt, telegram, sonoff diy (TBD), and yeelight (TBD).
- History storage which persist received messages in sqlite db
- Versatile mapping rules engine which defines how application should respond to received messages
- Actions executor which executes one or more actions in respond to received message
- REST API service to manage application: create rules, read device messages history, read registered devices catalogue
- Telegram as a channel for delivery various notifications and alerts and remote management
- Docker compose to deploy entire server, which consists of this [backend](https://github.com/fedulovivan/mhz19-go), frontend (TBD), [device-pinger](https://github.com/fedulovivan/device-pinger) service, [eclipse mosquitto](https://mosquitto.org/) message broker, [zigbee2mqtt](https://www.zigbee2mqtt.io/) zigbee bridge 

### Unified message structure

No matter which channel was used to receive a message, or which certain device has emitted that message, we pack every message into unified strusture for future handling

- Channel type - mqtt, telegram, sonoff, yeelight
- Device class - zigbee device, zigbee bridge, device-pinger, valve-manipulator, telegram-bot, sonoff device, yeelight device
- Device id - unique device identified specific for the certain device class. e.g. zigbee ieee device address 0x00158d0004244bda
- Payload - the message itself as a json. untyped, specific for the certain device and channel. e.g. zigbee wall switch message may look like `{"battery":100,"action":"single_left"}` or telegram-bot message as `{"Text":"/open-valves"}`
- Timestamp - a time when message was received by the server, usefull when reading and visualizing historical data

### Used technologies

- [golang](https://go.dev/) as main app language
- [sqlite3](https://www.sqlite.org/) with [go-sqlite3](github.com/mattn/go-sqlite3) client as a persistent storage
- [ozzo-routing](github.com/go-ozzo/ozzo-routing/v2) for lightweight http router implementation
- [godotenv](github.com/joho/godotenv) and [go-envconfig](github.com/sethvargo/go-envconfig) as configuration layer
- [eclipse paho](github.com/eclipse/paho.mqtt.golang) as mqtt client
- [docker](https://www.docker.com/) for containerization
- [Makefile](./blob/main/Makefile) for developer routine automation
- [telegram bot api](https://core.telegram.org/bots/api) with [client](https://github.com/go-telegram-bot-api/telegram-bot-api) for the notifications and remote management