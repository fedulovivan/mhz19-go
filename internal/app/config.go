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
	DbDebug        bool   `env:"DB_DEBUG,default=false"`
	SqliteFilename string `env:"SQLITE_FILENAME,default=database.bin"`

	// telegram
	TelegramDebug  bool   `env:"TELEGRAM_DEBUG,default=false"`
	TelegramToken  string `env:"TELEGRAM_TOKEN"`
	TelegramChatId int64  `env:"TELEGRAM_CHATID"`

	// mqtt
	MqttDebug    bool   `env:"MQTT_DEBUG,default=false"`
	MqttHost     string `env:"MQTT_HOST,default=mosquitto"`
	MqttPort     int    `env:"MQTT_PORT,default=1883"`
	MqttUsername string `env:"MQTT_USERNAME"`
	MqttPassword string `env:"MQTT_PASSWORD"`
	MqttClientId string `env:"MQTT_CLIENT_ID,default=mhz19-go"`

	// other
	Tz            string        `env:"TZ"`
	LogLevel      slog.Level    `env:"LOG_LEVEL,default=debug"`
	BuriedTimeout time.Duration `env:"BURIED_TIMEOUT,default=90m"`
	RestApiPort   int           `env:"REST_API_PORT,default=8888"`
}

func InitConfig() {
	RecordStartTime()
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
	configAsJson, _ := json.MarshalIndent(Config, "", "  ")
	fmt.Printf("starting with config %v\n", string(configAsJson))
	if Config.IsDev {
		fmt.Println("all known config variables", getExpectedEnvVars())
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
