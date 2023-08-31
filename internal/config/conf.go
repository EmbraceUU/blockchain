package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	APIServer APIServer `json:"apiServer"`
	BTC       BTC       `json:"btc"`
	EVM       EVM       `json:"evm"`
}

type APIServer struct {
	Port string `json:"port"` // 端口
}

type EVM struct {
	EndPoint string `json:"endPoint"`
}

type BTC struct {
	ID       string `json:"id"`
	EndPoint string `json:"endPoint"`
}

func NewConfig() Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("../../api/")
	v.AddConfigPath("./api/")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("init config failed, %s", err.Error())
	}

	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		log.Fatalf("init config failed, %s", err.Error())
	}

	return config
}
