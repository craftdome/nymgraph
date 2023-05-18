package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	configName string

	UseProxy      bool   `yaml:"UseProxy"`
	Proxy         string `yaml:"Proxy"`
	TestProxySite string `yaml:"TestProxySite"`

	NymClient struct {
		Host string `yaml:"Server"`
	} `yaml:"NymClient"`
}

func NewConfig(configName string) (*Config, error) {
	cfg := &Config{configName: configName}
	data, err := os.ReadFile(configName)
	if err != nil {
		return nil, errors.Wrap(err, "ReadFile")
	}

	// Парсим конфиг
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return cfg, nil
}

func (c *Config) Save() error {
	file, err := os.OpenFile(c.configName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrap(err, "OpenFile")
	}
	defer file.Close()

	if data, err := c.Bytes(); err != nil {
		return errors.Wrap(err, "Bytes")
	} else if _, err = file.Write(data); err != nil {
		return errors.Wrap(err, "Write")
	}

	return nil
}

func (c *Config) Bytes() ([]byte, error) {
	if data, err := yaml.Marshal(c); err != nil {
		return nil, errors.Wrap(err, "Marshal")
	} else {
		return data, nil
	}
}
