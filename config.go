package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MagicStringRule struct {
	Enabled   bool `yaml:"enabled"`
	MaxLength int  `yaml:"max-length"`
}

type ParamCountRule struct {
	Enabled bool `yaml:"enabled"`
	Max     int  `yaml:"max"`
}

type FuncLengthRule struct {
	Enabled  bool `yaml:"enabled"`
	MaxLines int  `yaml:"max-lines"`
}

type Config struct {
	Rules struct {
		CheckExportedComments bool            `yaml:"check-exported-comments"`
		CheckMagicStrings     MagicStringRule `yaml:"check-magic-strings"`
		CheckParameterCount   ParamCountRule  `yaml:"check-parameter-count"`
		CheckFunctionLength   FuncLengthRule  `yaml:"check-function-length"`
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
