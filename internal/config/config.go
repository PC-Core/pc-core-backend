package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	DbDriver string `yaml:"dbdriver"`
}

func ParseConfig(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
