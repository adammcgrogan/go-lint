package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MagicStringRule struct {
	Enabled   bool `yaml:"enabled"`
	MaxLength int  `yaml:"max-length"`
}

type Config struct {
	Rules struct {
		CheckExportedComments bool            `yaml:"check-exported-comments"`
		CheckMagicStrings     MagicStringRule `yaml:"check-magic-strings"`
	} `yaml:"rules"`
}

func loadConfig() (*Config, error) {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
