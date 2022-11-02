package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbUrl         		string `mapstructure:"DB_URL"`
	JWTString     		string `mapstructure:"SECURE_STRING"`
	ExternalUrl   		string `mapstructure:"EXTERNAL_URL"`
	Adress        		string `mapstructure:"ADRESS"`
	AllowedOrigns 		string `mapstructure:"ALLOWED_ORIGINS"`
	DbTestUrl 	  		string `mapstructure:"DB_TEST_URL"`
	AllowCredentials 	bool `mapstructure:"ALLOW_CREDENTIALS"`
}

func ReadConfig() (*Config, error) {
	config := new(Config)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./../../configs") // for tests
	viper.AddConfigPath("../configs") // for tests
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.Unmarshal(&config)
	return config, nil
}
