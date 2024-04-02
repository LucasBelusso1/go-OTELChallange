package config

import (
	"os"
)

var config Config

type Config struct {
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig() {
	config.WeatherApiKey = os.Getenv("WEATHER_API_KEY")
	if config.WeatherApiKey == "" {
		panic("WEATHER_API_KEY must be set")
	}
}

func GetConfig() Config {
	return config
}
