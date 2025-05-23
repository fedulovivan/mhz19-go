package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

var Config ConfigStorage

type ConfigStorage struct {
	// development mode
	IsDev bool `env:"DEV,default=false"`

	// db
	DbDebug bool `env:"DB_DEBUG,default=false"`

	SqliteFilename    string `env:"SQLITE_FILENAME,default=./sqlite/database.bin"`
	SqliteBusyTimeout int    `env:"SQLITE_BUSY_TIMEOUT,default=5000"`
	// SqliteMaxTxDuration time.Duration `env:"SQLITE_MAX_TX_DURATION,default=60s"`

	// telegram
	TelegramDebug         bool     `env:"TELEGRAM_DEBUG,default=false"`
	TelegramTokens        []string `env:"TELEGRAM_TOKENS"`
	TelegramChatId        int64    `env:"TELEGRAM_CHAT_ID"`
	TelegramDefaultOutBot string   `env:"TELEGRAM_DEFAULT_OUT_BOT"`

	// mqtt
	MqttDebug    bool   `env:"MQTT_DEBUG,default=false"`
	MqttHost     string `env:"MQTT_HOST,default=mosquitto"`
	MqttPort     int    `env:"MQTT_PORT,default=1883"`
	MqttUsername string `env:"MQTT_USERNAME"`
	MqttPassword string `env:"MQTT_PASSWORD"`
	MqttClientId string `env:"MQTT_CLIENT_ID,default=mhz19-go"`

	// other
	Tz                   string        `env:"TZ"`
	LogLevel             slog.Level    `env:"LOG_LEVEL,default=debug"`
	DefaultBuriedTimeout time.Duration `env:"BURIED_TIMEOUT,default=90m"`
	RestApiPort          int           `env:"REST_API_PORT,default=8888"`
	RestApiPath          string        `env:"REST_API_PATH,default=/api"`
	ArgsDebug            bool          `env:"ARGS_DEBUG,default=false"`
	RulesFetchingLimit   time.Duration `env:"RULES_FETCHING_LIMIT,default=10s"`
	DnssdDebug           bool          `env:"DNSSD_DEBUG,default=false"`
}

func InitConfig() {
	// RecordStartTime()
	fileName, withConf := os.LookupEnv("CONF")
	if !withConf {
		fileName = ".env"
	}
	err := godotenv.Load(fileName)
	if err != nil {
		fmt.Println("godotenv.Load()", err)
	} else {
		fmt.Printf("env variables were loaded from %v file\n", fileName)
	}
	if err := envconfig.Process(context.Background(), &Config); err != nil {
		panic("failed loading env variables into struct: " + err.Error())
	}
	var configAsJson []byte
	configAsJson, _ = json.Marshal(Config)
	fmt.Printf("starting with config %v\n", string(configAsJson))
	if Config.IsDev {
		varsAsJson, _ := json.Marshal(getExpectedEnvVars())
		fmt.Println("known config variables", string(varsAsJson))
	}
}

// Use reflection to extract known config vars from ConfigStorage
func getExpectedEnvVars() []string {
	typ := reflect.TypeOf(ConfigStorage{})
	var m []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tagValue := field.Tag.Get("env"); tagValue != "" {
			tt := strings.Split(tagValue, ",")
			if len(tt) > 0 {
				m = append(m, tt[0])
			}
		}
	}
	return m
}

func GetMqttBrokerUrl() string {
	return fmt.Sprintf("tcp://%s:%v", Config.MqttHost, Config.MqttPort)
}
