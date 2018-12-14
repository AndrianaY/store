package config

import (
	"github.com/AndrianaY/store/util/baseconfig"
	"github.com/spf13/viper"
)

type settingsKeys struct {
	BucketID        string `name:"BUCKET_ID"`
	AppPort         string `name:"APP_PORT"`
	AppHost         string `name:"APP_HOST"`
	DbUser          string `name:"DB_USER"`
	DbPassword      string `name:"DB_PASSWORD"`
	DbServer        string `name:"DB_SERVER"`
	DbPort          string `name:"DB_PORT"`
	DbSchema        string `name:"DB_SCHEMA"`
	GoogleProjectID string `name:"GOOGLE_PROJECT_ID"`
	// LoggerName      string `name:"LOGGER_NAME" default:"MMLogger"`
}

var Keys = settingsKeys{}

func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	baseconfig.ReadConfigValues(&Keys)
}
