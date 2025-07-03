package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Log       LogConfig       `env:",prefix=LOG_"`
	PrusaLink PrusaLinkConfig `env:",prefix=PRUSALINK_"`
	Mqtt      MqttConfig      `env:",prefix=MQTT_"`
}

type LogConfig struct {
	Level  string `env:"LEVEL,default=INFO"`
	Format string `env:"FORMAT,default=text"`
}

type PrusaLinkConfig struct {
	Host     string `env:"HOST,required"`
	ApiKey   string `env:"APIKEY,required"`
	Interval int    `env:"INTERVAL,default=10"`
}

type MqttConfig struct {
	Broker   string `env:"BROKER,required"`
	Port     int    `env:"PORT,default=1883"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Topic    string `env:"TOPIC,default=prusa-link"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
