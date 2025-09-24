package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

var defaultConfig = Config{
	Domain: "example.com",
	Port:   20722,
	Log: ConfigLog{
		Level:  "info",
		Pretty: true,
	},
	Database: ConfigDatabase{
		Type: "sqlite",
		Url:  ":memory:",
	},
}

type Config struct {
	Domain   string         `yaml:"domain"`
	Port     int            `yaml:"port"`
	Log      ConfigLog      `yaml:"log"`
	Database ConfigDatabase `yaml:"database"`
}

type ConfigLog struct {
	Level  string `yaml:"level"`
	Pretty bool   `yaml:"pretty"`
}

type ConfigDatabase struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

func ReadConfig(configPath string) (config Config, err error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, err
	}

	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return defaultConfig, err
	}
	return
}
