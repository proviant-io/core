package config

import (
	"gopkg.in/yaml.v3"
	"io"
)

type DB struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

type Config struct {
	Db DB `yaml:"db"`
	// db
	// api or full
	// image storage type & location
}

func NewConfig(r io.Reader) (*Config, error) {

	cfg := &Config{}

	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
