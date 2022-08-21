package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbUrl string
	JWTString string
}

func ReadConfig() (*Config, error) {
	var config *Config = &Config{}
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err 
	}
	config.DbUrl = viper.GetString("dburl")
	config.JWTString = viper.GetString("secure_string")
	return config, nil
}