package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mysql struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"mysql"`
	Redis struct {
		Ip   string `yaml:"ip"`
		Port string `yaml:"port"`
	}
	Key string `yaml:"key"`
}

func NewConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
