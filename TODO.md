### Prio 0
- (+) feat: calculate average for counters.Time
- (+) complete race tests for counters.Time

### Prio 1
- api: toggle rule on/off
- api: toggle device buried_timeout on/off
- api: update rule? looks its better to use delete/create strategy

### Bugs
- bug: find the reason of no sound on macmini
- bug: check why docker build always takes 203s on macmini (Building 202.9s (17/17) FINISHED)
- bug: no mqtt (re)connection if network was not available on app startup and returned online later
- bug: "http: superfluous response.WriteHeader call from github.com/go-ozzo/ozzo-routing/v2.(*Router).handleError (router.go:131)" - appears after interruption of progressing apache bench - need to ensure this is normal, and not an application-level issue
- bug: "apr_socket_recv: Operation timed out (60)" - https://stackoverflow.com/questions/30352725/why-is-my-hello-world-go-server-getting-crushed-by-apachebench, try to find protection
- bug: "api:getAll took 3.451973917s" when reading 1k rules 1k times - try same scenario with postgres - ensure there is no room for optimisation here

### Features
- feat: parse DeviceClass(telegram-bot) as well as DeviceClass(5)
- feat: parse DeviceClass(mqtt) as well as ChannelType(1)
- feat: add room entity, connect it with devices
- (+) feat: better api for counters.Time()
- feat: detect bot(s) are connected/started, instead of using dumb timeout before publishing "Application started" message - decoupling and introducing outgoing queue may help here
- feat: do auto db backup before running any kind of migration tasks
- feat: ability to disable certain condition or action
- feat: log rule/condition/action executions to the db table
- feat: simple frontend
- feat: introduce rules.comments column
- feat: merge Zigbee2MqttSetState and ValveSetState actions
- feat: create meta which descibes expected args for conditions and actions and validate in rest api

### Arch changes/decisions
- arch: introduce same approach for handling outgoing messages: action only submits new message to out channel, and corresponding provider handles it asynchronous manner + store outgoing messages history as well
- arch: align arg names across actions (like we have two spellings: Cmd and Command)
- arch: think how (where?) we can construct/init "TemplatePayload" automatically, now we need to build it manually in action implementation
- mile: device_id + device_class adressing issue:
    - "FOREIGN KEY (device_id) REFERENCES devices(native_id)" requires sole UNIQUE index for column devices.native_id, while we actually need UNIQUE(device_class_id, native_id) since its unreasonable to constraint native_id across devices off all classes
    - in addition to "native_id" problem see also "unsafemap" in internal/entities/ldm/repository.go
    - a solution could be to keep ids as strings like "ZigbeeDevice(0x00158d000a823bb0)" or "Pinger(192.168.88.44)"
    - Message, make fields DeviceClass and DeviceId optional
    - arch: think of good api (constructor) for creating new message (NewMessage) Id, Timestamp, ChannelType, DeviceClass?, DeviceId? -  are mandatory
- arch: in NewEngine create mocks for all services, which will panic with friendly message if user forgot to set that service
- arch: get rid of any in Send(...any) - no ideas so far
- arch: split rest api and engine into different microservices
- arch: consider replacing sql.NullInt32 and sql.NullString with corresponding of pointer types - https://stackoverflow.com/questions/40092155/difference-between-string-and-sql-nullstring, for now stick with existing approach as more convenient
- arch: switch to nil instead of sql.NullInt32 - easy to MarshalJSON

### Try
- try: to deploy on old rpi/raspberrypi with ram disk enabled
- try: hcl - https://github.com/hashicorp/hcl
- try: some interactive cli framework for provision tool like cobra
- try: validation https://github.com/asaskevich/govalidator OR https://github.com/go-ozzo/ozzo-validation
- try: find out why cli command "make test" and "vscode" report different coverage statistics: 86.9% vs 100%. vscode syntax - `Running tool: /opt/homebrew/bin/go test -timeout 30s -coverprofile=/var/folders/5v/0wjs9g1948ddpdqkgf1h31q80000gn/T/vscode-go7lC7ip/go-code-cover github.com/fedulovivan/mhz19-go/internal/engine`
- try: separate di library https://pkg.go.dev/go.uber.org/fx
- try: openapi or swagger https://en.wikipedia.org/wiki/OpenAPI_Specification or https://swagger.io/
- try: http router https://github.com/julienschmidt/httprouter istead of ozzo-routing
- try: prometheus https://prometheus.io/docs/guides/go-application/, https://habr.com/ru/articles/709204/
- try: grpc
- try: benchmarking tool https://github.com/sharkdp/hyperfine
- try: postgres instead of sqlite3
- try: mongodb instead of sqlite3
- try: chatgpt or copilot to review code https://www.reddit.com/r/vscode/comments/14upva0/how_to_use_chatgptcopilot_for_code_review/
- try: https://github.com/mheffner/go-simple-metrics, https://github.com/hashicorp/go-metrics or release own
- try: wrk utility (analog of ab, hey, oha) https://github.com/wg/wrk
- try: yandex-tank https://github.com/yandex/yandex-tank

### Milestones

- (+) 26 sep 2024, DONE. Initial launch - All features from mhz19-next plus storing mapping rules in database
- Implement simple frontend
- Interactive zigbee device join - End-to-end scenario with new device device join, confuguring rules, with no app retart
- device_id + device_class adressing issue

### Completed

- (+) bug: DATA RACE when strarting with incorrect mqtt settings
- (+) feat: think how to design "bot reply feature" - TelegramBotMessage has access to initial message from 'telegram-bot' channel which contains ChatId in payload
- (+) bug: `requesting message for otherDeviceId=192.168.88.44` Not Equal gives wrong result (added more logging, probably caused by recent app restart and missing data to execute otherDeviceId logic and get actual pinger status for 192.168.88.44) real reason was in using ldm service Get(m.DeviceClass, otherDeviceId) instead of GetByDeviceId(otherDeviceId)
- (+) bug: Msg=40 Rule=13 Condition=48 Changed Started args=map[Value:$message.action] res=false - regression
- (+) bug: "no last message for.." should not be recorded as ERROR
- (+) bug: last device message is not recorded if no rules are configured
- (+) api: delete device
- (+) try: compile with "race" flag - https://www.youtube.com/watch?v=mvUiw9ilqn8&list=WL&index=4
- (+) arch: use two telegram channels for notifications: for CRITICAL messages and all the rest.
- (+) arch: think how to distinquish "end device" message from all "others" - just as new flag for the message Struct?
- (+) arch: mapping rules could be pre-defined (system) and loaded from db (user-level) - think we need to store everything in db, even system rules
- (+) feat: min/max time metrics
- (+) api: rename device, add device, delete rule
- (+) bug: DevicesService.UpsertAll return ids
- (+) feat: avoid provisioning "devices" from 00-sql-up script, leave only system devices or get rid even from them
- (+) bug: a lots of rules created after load tests block application slow down app startup, with no details in logs - for now wrap in goroutine and added more logging
- (+) switched to oha for load testing, fixed "rules-write" scenario, eliminated "ab" and "hey"
- (+) try: https://github.com/hatoo/oha
- (+) bug: avoid "PostSonoffSwitchMessage Start args="map[Command:off DeviceId:10012db92b]"", check the status? use change?
- (+) bug: debug logging issues: log-level-1.txt, log-level-2.txt, log-level-3.txt, see new test internal/logger/tag_test.go - wrong approach for cloning structure with underlying slice in Tag::With()
- (+) bug: at present moment there is no "previous message info" in ExecuteActions. we call with ONE message for non-throttled rule, and with ARRAY of messages if throttling is enabled. so flag IsFirst is implemented and handled incorrectly
- (+) arch: introduce new table "schema version" and validate migration number during startup
- (+) bug: add tag to the args reader logger
- (+) bug: reader.Get("Value"): in=$message.occupancy (string), out=<nil> (<nil>) -> Msg=14 Rule=6 condition=False Fail -> cannot cast to bool - fixed error message, added extra condition = !Nil to avoid error
- (+) nice: do concurrent execution of processFile in cmd/provisioning/main.go
- (+) nice: switch from $deviceId to fn=DeviceId in all conditions
- (+) feat: implement new action to play alert
- (+) bug: lots of erorrs: `ERR [engine]     Msg=1906 Rule=6 condition=InList Fail err="Message.ExecDirective(): Payload 'map[string]interface {}, map[battery:100 device_temperature:30 linkquality:90 power_outage_count:24 voltage:3025]' has no field 'action'"` - for now fixed with supressing message "has no field", look redundant
- (+) feat: issue zigbee2mqtt/bridge/config/devices/get periodically - not supported by z2m now, zigbee2mqtt/bridge/devices is retained message, updated by z2m with list of devices
- (+) arch: avoid postman as a dependency for provisioning system/user rules. create them via curl script - introduced appropriate requests json in assets + new makefile command "provisioning"
- (+) arch: make BotName optional for the TelegramBotMessage, take it from config
- (+) feat: find a place for "Application started" message
- (+) arch: reworked counters module / per-rule match counter
- (+) arch: try benchmark tests
- (+) bug: fixed logging/execution issue for InvokeActionFunc
- (+) feat: creted all "legacy" rules
- (+) bug: use host network to fix multicast in docker - https://github.com/flungo-docker/avahi
- (+) feat: migrate all rules from mhz19-next
- (+) bug: avoid "%!,(MISSING)" in logs - caused by usage of tag.F with inner fmt.Sprintf
- (+) bug: no DeviceId prefix in serialized json - https://stackoverflow.com/questions/39164471/marshaljson-not-called
- (+) feat: create api to read device classes (or unified api for any simple dict table?)
- (+) bug: cannot use "otherDeviceId" as "DeviceId(192.168.88.44)" need "192.168.88.44"
- (+) queued message may be logged with initial message id - "23:30:12.390 DBG [engine] Msg#1004 Rule=3 action=RecordMessage End" - tag fn bound to msg and rule "captured" by queue flush callback
- (+) bug: db Tid is not unique within transaction
- (+) feat: enable throttling for "RecordMessage" action + implement batch insert + ensure there no misses with throttled handling
- (+) introduce new channel - rest
- (+) bug: second request to sonoff hangs - not closed body reader
- (+) `make docker-up` uses REST_API_PORT from Makefile, switch to .env - https://stackoverflow.com/questions/44628206/how-to-load-and-export-variables-from-an-env-file-in-makefile
- (+) arch: for rule_condition_or_action_arguments use value + data_type_id instead in addition to + device_id + device_class_id + channel_type_id
- (+) uts: create tests for nested conditions
- (+) arch: avoid "args=map[]" in logs - slog always writes nil map as "map[]", not "nil", see also Test20, Test21 in service_test.go, https://github.com/golang/go/issues/69496
- (+) bug: unable to start vscode debugging or unit tests with current implementation around SqliteFilename and SQLITE_FILENAME
- (+) arch: support several bots
- (+) try: gorm - /Users/ivanf/Desktop/Projects/go/gorm-test/main.go
- (+) feat: implement log tag with meta, so we can add attrs to function
- (+) feat: create api to read one device
- (+) feat: add devices.buried_ignored column or devices.buried_timeout (0 - blacklisted device, null - default timeout)
- (+) basic app counters
- (+) feat: add Condition.FnInverted bool flag instead of NotEqual, NotChannel
- (+) feat: introduce SkipCounter field for the rule
- (+) feat: create api to delete rule
- (+) bug: a weird sporadic "ERR [rest] Not Found" in logs for sussess responses - RCA this is chrome requesting /favicon.ico along with api url
- (+) arch: messages with ewelink device mdns announcement should not have device id in DeviceId field, since semantically this is not a message from device itself (same as "special" zigbee bridge message with devices list)
- (+) arch: rename system provider to "hnsfaw" or "notseen" or "buried"
- (+) feat: implement "buried devices" aka "have not seen for a while" notifications
- (+) feat: api: log errors captured by router error handler, also change default handler to render error as a json
- (+) bug: MarshalJSON is not working for condition.fn and throttle - change from pointer to value receiver
- (+) feat: add OtherDeviceId to repository, schema and service
- (+) big: fix "go-sqlite3 requires cgo to work" for docker build
- (+) feat: add Dockerfile
- (+) feat: finish implementation for "otherDeviceId"
- (+) feat: finish implementation of all actions
- (+) try: https://github.com/go-ozzo/ozzo-routing
- (+) bug: figure out why we cannot test engine in uts end to end - internal/engine/mappings_test.go::Test10, fixed with emitting message within timeout
- (+) quest: find out why Args::UnmarshalJSON() is not called in Test164 - just because go's encoding/json/decode.go::Unmarshal() contains call ofr checkValid(), which prevents further parsing if invalid input was given
- (+) arch: think how to move messages_* and devices_* outside of engine package to their own
- (+) feat: sonoff provider, mdns(dns-sd) client for sonoff devices https://github.com/hashicorp/mdns
- (+) feat: rename channel=sonoff to channel=mdns, add new rule for dnssd-sources devices insertion
- (+) bug: "database is locked" - dig deeper into https://github.com/mattn/go-sqlite3/issues/274
- (+) feat: $channelType directive
- (+) feat: load db mapping rules on engine startup
- (+) feat: implement last device messages api
- (+) bug: rule created via api is not loaded to engine
- (+) feat: implement template-based argument value mappings - https://pkg.go.dev/text/template
- (+) arch: extract actions in separate files
- (+) feat: implement throttle
- (+) feat: finish handling of "Mapping" in InvokeActionFunc (NewArgReader)
- (+) feat: add action.PayloadData property - PayloadData was replaced by action args in new design
- (+) bug: "s.repository.Get(db.NewNullInt32(ruleId))" returns wrong data - bug in db.AddWhere
- (+) internal/devices/repository.go::UpsertDevices should not stop on error - was actual for plain insert only
- (+) arch: get rid of full path in SQLITE_FILENAME to run tests - no more actual after introducing di
- (+) db file stuck on 6.5mb as with 12k even after migrate-reset - need a VACUUM after dropping tables
- (+) got an error "database is locked" executing INSERT INTO messages... - looks we need explicit tx.Rollback()
- (+) finish implementation of device.upsert
- (+) create devices api
- (+) create messages api
- (+) create action UpsertZigbeeDevices
- (+) delete internal/engine/dummy_provider.go - https://go.dev/blog/deadcode
- (+) insert unknown devices automatically
- (+) implement "record message" action
- (+) create messages service and make it depency of engine
- (+) implement /rules/1 endpoint and where querying
- (+) implement /stats endpoint
- (+) double logging for "DBG [stats]   ✨ repo:Get took 755.584µs" - just forgotten extra call of slog in TimeTrack
- (+) create table messages
- (+) develop an approach of passing device id "0x00158d00042446ec" via json (unmarshalling)
- (+) finish ToDbArguments
- (+) implement REST API to read/create/update rules (use https://github.com/go-ozzo/ozzo-routing or https://github.com/gin-gonic/gin)
- (+) handle DbRuleConditionOrActionArgument.IsList and create tests
- (+) create readme and license
- (+) reorganise code to conform service repository pattern (https://medium.com/@ankitpal181/service-repository-pattern-802540254019); internal/rest/rest.go > rules/api; internal/rest/service.go > rules/service; internal/db/model.go > rules/repository; move engine.BuildMappingRules, engine.FlattenConditions to service layer
- (+) finish BuildRules (BuildConditions etc)
- (+) implement Rule from/to Json layer
- (+) ON DELETE CASCADE is not working - for sqlite requires "PRAGMA foreign_keys=ON" to work
- (+) get rid of three "func init()" (use DI?) - for now just manuall call of Init(s) in main
- (+) "transaction has already been committed or rolled back" on load test - refuse from using BeginTx with ctx from errgroup
- (+) refactor FetchAll to concurrent call, use contexts
- (+) try https://github.com/jmoiron/sqlx, https://jmoiron.github.io/sqlx/ - library has lots of issues (300 open, 370 closed) and , for now using own lightweight wrappers
- (+) no new mqtt messages after mqtt disconnect/autoreconnect (`Connection lost error="pingresp not received, disconnecting"` and later `Connected broker=tcp://macmini:1883`) + same issue for device-pinger which impacts its service - subscriptions should be settled in connect handler
- (+) schema, use same approach to check fk for device_classes and function_name (either CHECK constraint of FK to separate table)
- (+) add transaction id to handleMessage log records
- (+) execute actions asyncronously in goroutine
- (+) "List" argument could be only defined as []any, not []string, not []int etc - no solution so far
- (+) think how to organize OutChannel validation for the actions which require it - use callback
- (+) add more debug
- (+) think how to revive Service.Type() which is looking redundant now
- (+) add support of hi-level composite functions like ZigbeeDevice under the hood composed of "standart" Equal and InList
- (+) Service.Receive() schould be placed in "base" struct - the only one working idea is to introduce separate folder service/base and corresponding separate package base_service with public base struct and public fields within. then import it to mqtt and tbot and embed.
- (+) refactor to avoid "startup messages" in unit tests - logger initialisation moved to options
- (+) wrap internal/mqtt/client.go and internal/tbot/tbot.go into structs
- (+) for mqtt client, rather than hardocding in defaultMessageHandler, define rules/adapters for transforming topic and payload into final message per device class
- (+) consider replacing hand-written adapters with mqttClient.AddRoute() API, also add warning for messages captured by defaultMessageHandler (assuming all topics we subscribe should have own handlers and default one should not be reached)

### Discarded
- (?) arch: make logger and logTag a dependency of service, api and repository - no urgent need, everything is easily testable with current approach
- (?) try: opentelemetry https://opentelemetry.io/docs/languages/go/getting-started/, https://www.reddit.com/r/devops/comments/nxrbqa/opentelemetry_is_great_but_why_is_it_so_bloody/ - no need now, due to overcomplicated api
- (?) introduce intermediate layer between named args and function implementation using regular args (more robust, simplify things like ZigbeeDeviceFn)
- (?) think about "first match" strategy in handleMessage - we do not need this, since we to execute RecordMessage and some other action 
- (?) feat: create test service for sonoff wifi devices (poll them periodically to receive status updates)
- (?) arch: looks like we need to compare values as srings in conditions

### new mapping rule structure
```golang
    type ConditionArgument interface {
        int | string | bool
    }
    type LogicOp string
    type ConditionFn string
    type ActionFn string
    const AND ConditionOp = "AND"
    const OR ConditionOp = "OR"
    type Action struct {
        Fn ActionFn
        Args []ConditionArgument
        Mapping
    }
    type Condition struct {
        Fn ConditionFn
        Args []ConditionArgument
        LogicOp LogicOp
        List []Condition
    }
    type Rule struct {
        Condition Condition
        Actions []Action
        Throttle int
    }    
```

// example 1
// execute RelayToTelegram action for all messages with deviceId=192.168.88.1
// (note that channel and device class are ignored)
{
    Condition{ Fn: "Equals" Args: ["$deviceId", "192.168.88.1"] }
    Actions{ RelayToTelegram }
}

// example 2
// execute RelayToTelegram, SaveToSqlite actions for all messages
// which conform
// deviceClass=ZIGBEE AND deviceId=0x00158d0000c2fa6e AND (Changed message.state OR NotNil message.lastSeen)
{
    Condition{Op AND, List{
        { Fn: "Equals" Args: ["$deviceClass", ZIGBEE] }
        { Fn: "Equals" Args: ["$deviceId", "0x00158d0000c2fa6e"] }
        { Op: OR List{
            { Fn: "Changed" Args: ["$message.state"] }
            { Fn: "NotNil" Args: ["$message.lastSeen"] }
        } }
    }}
    Actions{ RelayToTelegram, SaveToSqlite }
}

// example 3
// implement high-level helper functions
// func ZigbeeDevice () Equals deviceClass AND  
{
    Condition{ Fn: "ZigbeeDevice" Args: ["0x00158d0000c2fa6e", "0x00158d000405811b"] }   
}

// rules
//   id int
//   comments string
//   enabled bool
//   throttle int
// rule_conditions
//   id
//   rule_id
//   function_name
//   logical_operation
//   parent_condition_id
// rule_actions
//   id
//   rule_id
//   function_name
//   device_id
// rule_condition_and_action_arguments
//   id
//   rule_id
//   condition_id
//   action_id
//   string_value
//   device_id_value
// rule_action_mappings
//   id
//   rule_id
//   action_id
//   key
//   value