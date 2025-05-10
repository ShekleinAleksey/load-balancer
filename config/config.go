package config

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ListenPort  string      `mapstructure:"listen_port"`
	Algorithm   string      `mapstructure:"algorithm"`
	Backends    []string    `mapstructure:"backends"`
	RateLimit   RateLimit   `mapstructure:"rate_limit"`
	HealthCheck HealthCheck `mapstructure:"health_check"`
}

type HealthCheck struct {
	Interval string `mapstructure:"interval"`
	Timeout  string `mapstructure:"timeout"`
}

type RateLimit struct {
	Capacity     int           `mapstructure:"capacity"`
	RefillPerSec time.Duration `mapstructure:"refill_per_sec"`
}

func InitConfig() (config Config, err error) {
	viper.AddConfigPath("config")
	viper.AddConfigPath("/app/config") // Для Docker
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("ошибка чтения конфига: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	if config.ListenPort == "" {
		config.ListenPort = "8080" // значение по умолчанию
	}

	logrus.Info("Конфиг: ", spew.Sdump(config))

	return config, nil
}
