package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPListen string `yaml:"http_listen"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}

func loadConfig(configFile string) (config *Config, err error) {
	configContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config = &Config{}
	err = yaml.Unmarshal(configContent, config)
	if err != nil {
		return nil, err
	}

	return
}
