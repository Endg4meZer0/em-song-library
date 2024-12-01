package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port           int    `mapstructure:"PORT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	ExternalAPIURL string `mapstructure:"EXTERNAL_API_URL"`
}

func Load() (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(".env file is not found: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("An error occured during config unmarshalling: ", err)
	}

	return config, nil
}
