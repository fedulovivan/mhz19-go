### Prio 0
none

### Prio 1
none

### Bugs
- bug: avoid MQTT_HOST=192.168.88.18 and MQTT_HOST=192.168.88.18 in .env config
- bug: docker: find the reason of rebuilding mhz19-frontend along with rebuild of mhz19-go; neither --always-recreate-deps, --no-deps helps nor removing depends_on helps; see https://github.com/docker/compose/issues/9600
- bug: macmini: find the reason of no sound (not actual after migration to rpi5)
- bug: reply to /door command send in chat with @Mhz19AlertsBot bot is send to @Mhz19Bot (default one from config)

### Digs
- dig: check why "api:getAll took 3.451973917s" when reading 1k rules 1k times; try same scenario with postgres; check there is no room for optimisation here
- dig: check why HandleMessage gives 1500 nanoseconds in benchmark, while prometheus measure is x1000 = 1.5-3milliseconds
- dig: read more about makefile PHONY https://vsupalov.com/makefile-phony/
  
### Features
- feat: deploy grafana with dashboards
- feat: api: add "?window=1h" to /messages/device/<deviceId>
- feat: uts: create unit tests for internal/providers/buried_devices/provider.go
- feat: instrument queue and queue container modules
- feat: sql: avoid ON DELETE CASCADE for the columns dependand on dictionaries (e.g to avoid unexpected loss of rules after reducing dictionary with actions list)
- feat: move docker-compose stack-related items into the separate repository; tidy up all volume-targeted folders/files structure with 3pp configs (prometheus, zigbee2mqtt, mosquitto etc), avoid anonymous volume for zigbee2mqtt; pass all zigbee2mqtt settings (host, serial/port, frontend port) via env vars and remove zigbee2mqtt-data/configuration.yaml from vcs
- feat: think how to init SqliteMaxTxDuration in unit tests, now app.InitConfig is not called in UTs
- feat: collect "messages by device id" metric
- feat: db: introduce updated_at, created_at columns
- feat: db: limit rule name length, since its used to group prometheus metrics
- feat: api: toggle rule on/off
- feat: api: add url /devices/id/15
- feat: api: toggle device buried_timeout on/off
- feat: api: update rule - looks its better to utilize delete/create strategy
- feat: accept DeviceClass(telegram-bot) as well as DeviceClass(5) in json
- feat: accept DeviceClass(mqtt) as well as ChannelType(1) in json
- feat: db: add new table "rooms", connect it with devices
- feat: db: auto db backup before running any kind of migration tasks
- feat: db: ability to disable certain condition or action
- feat: log rule/condition/action execution history to separate db table
- feat: db: introduce rules.comments column
- feat: create meta which descibes expected args for conditions and actions and validate them

### Arch changes/decisions
- arch: introduce same approach for handling outgoing messages: action only submits new message to out channel, and corresponding provider handles it asynchronous manner + store outgoing messages history as well
- arch: align arg names across actions (like we have two spellings: Cmd and Command)
- arch: think how (where?) we can construct/init "TemplatePayload" automatically, now we need to build it manually in action implementation
- arch: in NewEngine create mocks for all services, which will panic with friendly message if user forgot to set that service
- arch: get rid of any in Send(...any) - no ideas so far
- arch: split rest api and engine into different microservices
- arch: consider replacing sql.NullInt32 and sql.NullString with corresponding of pointer types - https://stackoverflow.com/questions/40092155/difference-between-string-and-sql-nullstring, for now stick with existing approach as more convenient
- arch: switch to nil instead of sql.NullInt32 - easy to MarshalJSON
- arch: device_id + device_class adressing issue (see more detailed tasks breakdown in "Milestones" section below)

### Try
- try: victoria logs instead of dozzle
- try: swagger/swaggo or something similar https://www.reddit.com/r/golang/comments/180jgzi/how_do_you_provide_documentation_for_your_rest/
- try: openapi or swagger https://en.wikipedia.org/wiki/OpenAPI_Specification or https://swagger.io/
- try: https://github.com/VictoriaMetrics/metrics instead of prometheus lib
- try: mongodb instead of sqlite3
- try: go version manager https://github.com/moovweb/gvm
- try: to utilize tcpdump to capture dnssd messages
- try: lib for online deadlock detection in go https://github.com/sasha-s/go-deadlock
- try: custom firmware for roborock vacuum https://valetudo.cloud/, https://www.youtube.com/watch?v=r_04K5SPEXI
- try: to disable go telemetry (/root/.config/go/telemetry/local)
- try: validation: https://github.com/go-playground/validator OR https://github.com/asaskevich/govalidator OR https://github.com/go-ozzo/ozzo-validation
- try: grpc
- try: to deploy on old/oldest rpi/raspberrypi with read-only fs enabled
- try: HashiCorp configuration language https://github.com/hashicorp/hcl
- try: some interactive cli framework for provisioning tool like cobra
- try: find out why cli command "make test" and "vscode" report different coverage statistics: 86.9% vs 100%. vscode syntax - `Running tool: /opt/homebrew/bin/go test -timeout 30s -coverprofile=/var/folders/5v/0wjs9g1948ddpdqkgf1h31q80000gn/T/vscode-go7lC7ip/go-code-cover github.com/fedulovivan/mhz19-go/internal/engine`
- try: 3pp di library https://pkg.go.dev/go.uber.org/fx
- try: http router https://github.com/julienschmidt/httprouter istead of ozzo-routing
- try: benchmarking tool https://github.com/sharkdp/hyperfine
- try: postgres instead of sqlite3
- try: chatgpt or copilot to review code https://www.reddit.com/r/vscode/comments/14upva0/how_to_use_chatgptcopilot_for_code_review/
- try: wrk utility (analog of ab, hey, oha) https://github.com/wg/wrk
- try: yandex-tank https://github.com/yandex/yandex-tank
- try: once.Do instead of "singleton" pattern - https://blog.matthiasbruns.com/golang-singleton-pattern
- try: errors library https://github.com/ansel1/merry
- try: to fix https://github.com/fedulovivan/effective-waffle/issues/9
- try: mockery https://github.com/vektra/mockery https://www.youtube.com/watch?v=eYHCCht8eX4
- try: go generate
- try: create load test for mqtt channel with `mosquitto_pub`
- try: In-memory SQL engine in Go sql/driver for testing purpose https://github.com/proullon/ramsql (from https://youtu.be/UfeZ-bPFs10?si=3FZTWpvjNvqh3X24&t=217)
- try: to switch from docker and docker-compose to kubernetes + helm OR openshift
- try: to learn how GOMAXPROC and docker --proc are related
- try: to create client for miio devices udp port 54321 (yeelight smart ceiling light, robot vacuum), for now stuck with token fetching issue. links: https://github.com/aholstenson/miio, https://github.com/OpenMiHome/mihome-binary-protocol, https://github.com/maxinminax/node-mihome, https://github.com/nickw444/miio-go, https://github.com/marcelrv/XiaomiRobotVacuumProtocol, https://github.com/vkorn/go-miio, https://www.youtube.com/watch?v=m11qbkgOz5o

### Milestones
- (+) 15 march 2025, DONE. Moved to raspberrypi 5 ("mhz19-go migration to new host" in notes)
- (+) 26 sep 2024, DONE. Initial launch - All features from mhz19-next plus storing mapping rules in database
- Implement simple frontend
- Prepare for the public usage (real use cases)
- Interactive zigbee device join (pairing/interview/adding/joining) - End-to-end scenario with new device device join, confuguring rules, with no app retart - https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html#zigbee2mqtt-bridge-event
- device_id + device_class adressing issue
    - "FOREIGN KEY (device_id) REFERENCES devices(native_id)" requires sole UNIQUE index for column devices.native_id, while we actually need UNIQUE(device_class_id, native_id) since its unreasonable to constraint native_id across devices off all classes
    - in addition to "native_id" problem see also "unsafemap" in internal/entities/ldm/repository.go
    - a solution could be to keep ids as strings like "ZigbeeDevice(0x00158d000a823bb0)" or "Pinger(192.168.88.44)"
    - Message, make fields DeviceClass and DeviceId optional
    - arch: think of good api (constructor) for creating new message (NewMessage) Id, Timestamp, ChannelType, DeviceClass?, DeviceId? -  are mandatory

### Completed
- (+) feat: freeze image version for 3pp docker deps (zigbee2mqtt 1.36.1 commit: ffc2ff1, mosquitto etc)
- (+) bug major: get rid of mqtt host ip in .env file
- (+) bug: "panic: send on closed channel" during gracefull shutdown; /Users/ivanf/Desktop/panic-001.txt
- (+) feat: handle seeding single asset file (instead of whole dir); quick implementation is introducing a FILTER parameter
- (+) feat: write exact time like "5 mins ago" instead of "for a while"
- (+) feat: rename provisioning to seed
- (+) feat: add ForceFlush to avoid "Waiting for the 9 message queues to stop"
- (+) feat: add a sibling for "Have not seen" message, which will notify device is back online
- (+) feat: reload change in devices.buried_timeout on the fly; already supported
- (+) feat: refactor and simplify parse_base method
- (+) bug: app still continue to receive messages after receiving shutdown
- (+) bug: two forms: Have not seen «ed6af05f0d59» for a while; Have not seen «Device of class valves-manipulator, with id 18225» for a while; see execTemplate - this is expected
- (+) feat: device white/black list for RecordMessage action; corrected on espresence side, added whitelist
- (+) bug: rulesService.OnCreated() returns rule conditions and actions without real db ids which impacts the log, see ~/Desktop/rules-git-diff.diff; we can use "Build" to construct types.Rule back from db objects, but repository.Create should update db objects with real ids - fixed with simplest approach: re-fetching from db
- (+) bug: wrong id (appeared id from messages table while was expected from devices) "Msg=25 Rule=1 Action=1 UpsertZigbeeDevices Created id=72852" - apparently this is https://github.com/mattn/go-sqlite3/issues/30, however my case is quite complex to reproduce: UpsertZigbeeDevices performs bulk UPSERT in transaction 1; RecordMessage performs bulk UPSERT into devices and messages in transaction 2; closing, no new ideas, no go-related specific
- (+) feat: switch to docker-compose in `make update`
- (+) feat: detect bot(s) are connected/started, instead of using dumb timeout before publishing "Application started" message - decoupling and introducing outgoing queue may help here
- (+) feat: merge Zigbee2MqttSetState and ValveSetState actions
- dig: mqtt: 0xe0798dfffed39ed1 (sonoff zigbee relay do not present in retained mqtt message with devices)
- (+) bug: no warning when app stops with non-empty queues, implement queues Wait logic
- (+) bug: no mqtt (re)connection if network was not available on app startup and returned online later: Firstly was written custom initialConnectWithRetry function, then found "opts.SetConnectRetry(true)" option for paho
- (+) bug: use host network to fix multicast in docker - https://github.com/flungo-docker/avahi; https://github.com/moby/libnetwork/issues/2397; https://forums.docker.com/t/multicast-support-on-docker-containers/144255: multicast cannot go through network bridge boundary, so for docker we need to use either host network or run dnssd part of the app in separate container (also with host networking enabled) 
- (+) bug: execute db migrations also from within docker container, otherwize we depend on sqlite3 binary installed on the host system, more specifically problem is macmini has version 3.31.1, while mbp 3.39.5. as a result a new feature "DROP COLUMN" is not working on macmini (introduced in sqlite 3.35, also see https://github.com/mattn/go-sqlite3/issues/927)
- (+) try: gdu fast disk usage analyzer https://github.com/dundee/gdu
- (+) dig: check why docker build always takes 203s on macmini, even after no changes in go files (Building 202.9s (17/17) FINISHED); RCA is rebuilding from scratch every time, without involing go build cache; reading: https://dev.to/jacktt/20x-faster-golang-docker-builds-289n; https://docs.docker.com/build/cache/optimize/
- (+) feat: configure builds with compose, now we have to build all images manually with separate tasks (mhz19-go, device-pinger, mhz19-front)
- (+) separate folder for sqlite files
- (+) bug: "apr_socket_recv: Operation timed out (60)" - https://stackoverflow.com/questions/30352725/why-is-my-hello-world-go-server-getting-crushed-by-apachebench; RCA this is in ab and macos limitations, no need to handle in app
- (+) bug: "http: superfluous response.WriteHeader call from github.com/go-ozzo/ozzo-routing/v2.(*Context).WriteWithStatus (context.go:178)" appears after interruption of progressing load test; need to ensure this is expected and not an application-level issue; reprodution is invoking `wget http://localhost:7070/api/rules` and immediate `Ctrl+C` when db contains 20k rules; RCA: first we start to write response normally with 200 code, meaning WriteHeader is already called, then after client disconnect an error "write: broken pipe" is raised and handled by errorHandler which calls WriteHeader again on attenmpt to "push" a json with error details and 500 code. Additinal read for "broken pipe" https://stackoverflow.com/questions/43189375/why-is-golang-http-server-failing-with-broken-pipe-when-response-exceeds-8kb, https://medium.com/trendyol-tech/golang-what-is-broken-pipe-error-tcp-http-connections-and-pools-3988b79f28e5
- (+) bug: race: /Users/ivanf/Desktop/race000 - appeared after recent moving reading map out of critical section in WithTag
- (+) bug: no abortion of "Fetching rules is still running" on Ctrl+C - code was synchronously blocked by long running rules_service.Build
- (+) bug: perf: rules_service.Build takes crazy amount of time to transform db records into to the rules representation - O(n^5) complexity caused by lots of repeating inner samber/lo calls, refactorred to advance indexing utilized objects. processing time reduced from 300s to 0.1s
- (+) bug: get rid of github.com/samber/lo
- (+) bug: add logging of the key, at the place where message was initially queued: "Rule=4 message queue is flushed now key=zigbee-device-0x00158d00067cb0c9-Rule4 mm=3"
- (+) ensure we accept interfaces and return concrete types (structs)
- (+) feat: include draft mhz19-front to the main compose stack
- (+) bug: "slice bounds out of range"  when reading single rule, RCA is https://github.com/goccy/go-json/issues/526
- (+) bug: still lots of err="got an error \"database is locked\" after eliminating SetMaxOpenConns(1) - https://stackoverflow.com/a/35805826/1012298, switched back to SetMaxOpenConns(1) and enabled WAL
- (+) feat: move services "device-pinger", "zigbee2mqtt", "mosquitto" from mhz19-next to local compose file + dont forget to extend .env file with required variables
- (+) bug: check why lots of records have duplicates http://macmini:7070/api/messages/device/0x00158d00067cb0c9?tocsv=1, for now stick with no changes on backend side, originally all messages are emitted by zigbee device (see assets/duplicated-messages-log.txt)
- (+) bug: make with no args invokes docker build
- (+) feat: switch from Seq to atomics
- (+) try: perf: can we speed up HandleMessage? create benchmarks for HandleMessage and conditions: baseline performance is about 3700ns per operation (1 Equal condition); 60% performance gain in Tag module after switching to strings.Builder and reduced heap memory allocations making strings.Builder no-pointer value; 90% performance gain with excluded conditions, TimeTrack, Tag module, logging, prometheus
- (+) bug: only 50rps for api-load-push-message-write - there was a rps limit, renamed makefile commands to avoid future confusions
- (+) bug: reset error counter on app restart (already works. why?) - grafana shows LAST metric on most of widgets, so its expected to see zeroing after app restart
- (+) bug: perf: check why api time is x3 of sql call: Tx#13 Transaction took 5.621441ms -> api:getByDeviceId took 15.907499ms - lots of time spent on json encoding, also BuildMessages did not used advance slice allocation
- (+) bug: now message for "guarded doors were opened/closed when i'm not at home" when rule is triggered for the first time after app restart - not a bug, next time I've forgot about throttled messages logic
- (+) feat: configure both telegram channels
- (+) bug: critical: panic: interface conversion: interface {} is nil, not string (when pairing new device), most prpbably from internal/engine/actions/upsert_zigbee_devices.go
- (+) feat: show app is up/down on dashboard - up{instance="host.docker.internal:7070"}
- (+) bug: memory does not realesed after calling /api/messages (growth from 5 tj 20mb) - ok, was released later
- (+) make all channels unbuffered
- (+) try: prometheus
- (+) feat: calculate average for counters.Time
- (+) complete race tests for counters.Time
- (+) feat: better api for counters.Time()
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
- (+) add transaction id to HandleMessage log records
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

### Discarded / doubtful
- (?) arch: make logger and logTag a dependency of service, api and repository - no urgent need, everything is easily testable with current approach
- (?) try: opentelemetry https://opentelemetry.io/docs/languages/go/getting-started/, https://www.reddit.com/r/devops/comments/nxrbqa/opentelemetry_is_great_but_why_is_it_so_bloody/ - no need now, due to overcomplicated api
- (?) introduce intermediate layer between named args and function implementation using regular args (more robust, simplify things like ZigbeeDeviceFn)
- (?) think about "first match" strategy in HandleMessage - we do not need this, since we to execute RecordMessage and some other action 
- (?) feat: create test service for sonoff wifi devices (poll them periodically to receive status updates)
- (?) arch: looks like we need to compare values as srings in conditions
- (?) feat: collect device up/down metrics
- (?) try: visualize buried devices in grafana
- (?) bug: some sql alterations from migrations cannot be undone by rollback https://stackoverflow.com/questions/4692690/is-it-possible-to-roll-back-create-table-and-alter-table-statements-in-major-sql/56738277; in sqlite DDL is transactional
- (?) try: https://github.com/mheffner/go-simple-metrics, https://github.com/hashicorp/go-metrics or release own
- (?) try: to eliminate mutexes from all race-condition-prune sections in favour of atomics https://betterprogramming.pub/atomic-pointers-in-go-1-19-cad312f82d5b; conclusion: for now mutexes are used solely to protect maps, so there is no better alternatives aside of Mutex (no way to use atomics, sync.Map ist too usefull for learning purposes, RCU approach is quite clean and to much complex to be placed instead of Mutexes)
 
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

<!-- https://github.com/mattn/go-sqlite3/issues/30 -->
Reproduced once for me. No exect scenario, unfortunately, just to record here. Lib version **v1.14.22**.
In [my case](https://github.com/fedulovivan/mhz19-go/blob/5714115b99c62eb4b0d9471d6c605b5ec8ac9e8b/internal/entities/devices/repository.go#L94C6-L94C17) insert into two tables (devices, messages) were performing each in its own transaction. And LastInsertId returned an Id from the other table.

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