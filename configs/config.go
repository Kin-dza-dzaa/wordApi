package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbUrl     string
	JWTString string
	Password  string
	ExternalUrl string
}

func ReadConfig() (*Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	config.DbUrl = viper.GetString("dburl")
	config.JWTString = viper.GetString("secure_string")
	config.Password = viper.GetString("password_string")
	config.ExternalUrl = viper.GetString("external_url")
	return &config, nil
}
