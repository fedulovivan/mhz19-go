
### Prio 0
- feat: finish implementation of all actions
- feat: api: log errors captured by router error handler, also change default handler to render error as a json
- feat: create api to update/delete rules
- feat: create api to add/update/delete devices
- feat: sonoff provider, mdns client for sonoff devices https://github.com/hashicorp/mdns

### Prio 1
- figure out why we cannot test engine in uts end to end - internal/engine/mappings_test.go::Test10
- bug: "http: superfluous response.WriteHeader call from github.com/go-ozzo/ozzo-routing/v2.(*Router).handleError (router.go:131)" - appears after termination of stucked apache bench
- bug: find why UnmarshalJSON is not called in Test164
- bug: "apr_socket_recv: Operation timed out (60)" - https://stackoverflow.com/questions/30352725/why-is-my-hello-world-go-server-getting-crushed-by-apachebench
- bug: "🧨 api:getAll took 3.451973917s" when reading 1k rules 1k times - try with postgres
- arch: "FOREIGN KEY (device_id) REFERENCES devices(native_id)" requires sole UNIQUE index for column devices.native_id, while we actually need UNIQUE(device_class_id, native_id) since its unreasonable to contraint native_id across devices off all classes
- arch: in addition to "native_id" problem see also "unsafemap" in internal/last_device_message/repository.go

### Prio 2
- arch: think how to move messages_* and devices_* outside of engine package to their own
- feat: implement log tag with meta, so we can add attrs to function
- bug: no mqtt (re)connection if network was not available on app startup and returned online later
- feat: create meta which descibes expected args for conditions and actions and validate
- arch: get rid of any in Send(...any)
- ut: create tests for recursive conditions
- arch: mappings rules could be pre-defined (system) and loaded from db (user-level) - think we need to store everything in db, even system
- feat: create test service for sonoff wifi devices (poll them periodically to receive status updates)
- try: find out why cli command "make test" and "vscode" report different coverage statistics: 86.9% vs 100%. vscode syntax - `Running tool: /opt/homebrew/bin/go test -timeout 30s -coverprofile=/var/folders/5v/0wjs9g1948ddpdqkgf1h31q80000gn/T/vscode-go7lC7ip/go-code-cover github.com/fedulovivan/mhz19-go/internal/engine`
- try: validation https://github.com/asaskevich/govalidator
- try: postgres instead of sqlite3
- try: https://pkg.go.dev/go.uber.org/fx
- try: opentelemetry - https://opentelemetry.io/docs/languages/go/getting-started/   
- try: prometheus
- try: grpc
- try: openapi or swagger https://en.wikipedia.org/wiki/OpenAPI_Specification or https://swagger.io/
- try: https://github.com/go-ozzo/ozzo-routing
- arch: make logger and logTag a dependencies for service, api and repository

### Completed

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
- (+) no new mqtt messages after mqtt disconnect/autoreconnect (`Connection lost error="pingresp not received, disconnecting"` and later `Connected broker=tcp://macmini:1883`) + same issue for device-pinger which impacts its service - subscribtions should be settled in connect handler
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
- (?) introduce intermediate layer between named args and function implementation using regular args (more robust, simplify things like ZigbeeDeviceFn)
- (?) think about "first match" strategy in handleMessage - we do not need this, since we to execute RecordMessage and some other action
- (?) split rest api and engine into separate microservices - does no look much reasonable, since both are heavily rely on db layer


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