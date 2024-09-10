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

### Applied best practices

- Source code organisation folows SOLID/clean-architecture recommendations
- DI to enable better UTs coverage
- Folders layout in accordance with [project-layout](https://github.com/golang-standards/project-layout)
- With unit tests, race tests and code coverage stats enabled
- Sql schema and basic data migrations support
- Load tests for the REST API
- Makefile for the common developer tasks 
- Backlog with TODOs

### Architecture

- Channel providers which collect messages from different channels and devices into unified stream. Supported channels are mqtt, telegram, dns-sd, sonoff (TBD) and yeelight (TBD).
- History storage which persist received messages in sqlite db
- Versatile mapping rules Engine which defines how application should respond to received messages
- Actions executor which executes one or more actions in respond to received message
- REST API layer to manage application: create rules, read device messages history, read registered devices
- Telegram as a channel to delivery various notifications and alerts and remote controll
- Docker compose to deploy entire server, which consists of this [backend](https://github.com/fedulovivan/mhz19-go), frontend (TBD), [device-pinger](https://github.com/fedulovivan/device-pinger) service, [eclipse mosquitto](https://mosquitto.org/) message broker, [zigbee2mqtt](https://www.zigbee2mqtt.io/) zigbee bridge

### Unified message structure

No matter which channel was used to receive a message, or which certain device has emitted that message, we pack every message into unified strusture to be handled by Engine.

- Channel type - mqtt, telegram, dns-sd, sonoff, yeelight
- Device class - zigbee device, zigbee bridge, device-pinger, valve-manipulator, telegram-bot, sonoff diy-plug device, yeelight device
- Device id - unique device identifier, specific for the certain device class. e.g. zigbee ieee device address 0x00158d0004244bda or device IP
- Payload - the message itself as a json. untyped, specific for the certain device and channel. e.g. zigbee wall switch message may look like `{"battery":100,"action":"single_left"}` or telegram-bot message as `{"Text":"/open-valves"}`
- Timestamp - a time when message was received by the server, usefull when reading and visualizing historical data

### Used technologies

- [golang](https://go.dev/) as main app language
- [sqlite3](https://www.sqlite.org/) with [go-sqlite3](github.com/mattn/go-sqlite3) client as a persistent storage
- [ozzo-routing](github.com/go-ozzo/ozzo-routing/v2) http routing
- [godotenv](github.com/joho/godotenv) and [go-envconfig](github.com/sethvargo/go-envconfig) as configuration layer
- [eclipse paho](github.com/eclipse/paho.mqtt.golang) as mqtt client
- [docker](https://www.docker.com/) for containerization
- [Makefile](./blob/main/Makefile) for developer routine automation
- [telegram bot api](https://core.telegram.org/bots/api) with [client](https://github.com/go-telegram-bot-api/telegram-bot-api) for the notifications and remote management
- [dnssd](https://github.com/brutella/dnssd) as mdns client (sonoff smart devices discovery)
- [gabs](https://github.com/Jeffail/gabs) as json querier

### Demo

![console.png](assets/demo-01.png)

### Tools required for the development on the bare host

- go `wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz`, `rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz`, `export PATH=$PATH:/usr/local/go/bin`
- golangci-lint https://golangci-lint.run/welcome/install `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.3`
- gcc `sudo apt install gcc`
- delve `go install -v github.com/go-delve/delve/cmd/dlv@latest`
- sqlite3 `sudo apt-get install sqlite3`