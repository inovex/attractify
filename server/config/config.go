package config

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Environment string `json:"environment"`
	AuthKeyHex  string `json:"authKey"`
	AuthKey     []byte `json:"-"`
	Server      struct {
		Bind      string `json:"bind"`
		DebugMode bool   `json:"debugMode"`
		Cert      string `json:"cert"`
		Key       string `json:"key"`
	}
	DB        string `json:"db"`
	Analytics string `json:"analytics"`
	Stream    struct {
		Brokers []string `json:"brokers"`
		Topic   string   `json:"topic"`
	} `json:"stream"`
	SMTP struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		FromAddress string `json:"fromAddress"`
	} `json:"smtp"`
}

func Parse(file string) (*Config, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	config.AuthKey, err = hex.DecodeString(string(config.AuthKeyHex))

	return &config, err
}
