package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func Init() {
	config = viper.New()
	config.SetEnvPrefix("APP")
	config.AutomaticEnv()

	config.SetConfigName("loadcentral-admin")
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s\n", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
