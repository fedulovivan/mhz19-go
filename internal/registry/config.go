package registry

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	_ "time/tzdata"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

var Config ConfigStorage

type ConfigStorage struct {
	TelegramToken  string     `env:"TELEGRAM_TOKEN"`
	TelegramChatId int64      `env:"TELEGRAM_CHATID"`
	MqttHost       string     `env:"MQTT_HOST,default=mosquitto"`
	MqttPort       int        `env:"MQTT_PORT,default=1883"`
	MqttUsername   string     `env:"MQTT_USERNAME"`
	MqttPassword   string     `env:"MQTT_PASSWORD"`
	MqttClientId   string     `env:"MQTT_CLIENT_ID,default=mhz19-go"`
	LogLevel       slog.Level `env:"LOG_LEVEL,default=debug"`
	IsDev          bool       `env:"DEV,default=false"`
	Tz             string     `env:"TZ"`
}

// Use reflection to extract known config vars from ConfigStorage
func GetExpectedEnvVars() []string {
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

func GetMqttBroker() string {
	return fmt.Sprintf("tcp://%s:%v", Config.MqttHost, Config.MqttPort)
}

func init() {
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
	fmt.Printf("starting with config %+v\n", Config)
	if Config.IsDev {
		fmt.Println("all known config variables", GetExpectedEnvVars())
	}
}
