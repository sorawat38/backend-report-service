package config

import (
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway/config"
	"github.com/spf13/viper"
)

type Config struct {
	App        App
	Log        Log
	Gateway    Gateway
	AWSSession AWSSession
}

type App struct {
	Name string
	Port string
}

type Log struct {
	Env string
}

type Gateway struct {
	PaymentService config.PaymentServiceCfg
	MenuService    config.MenuServiceCfg
}

type AWSSession struct {
	Id     string
	Secret string
}

func InitConfig() (Config, error) {

	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("config/") // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, err
		}
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
