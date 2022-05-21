package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadDotEnv() {
	viper.AddConfigPath(".")

	viper.SetConfigName(".env")

	viper.SetConfigType("env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("error in load dotEnv: %v", err)
	}
}
