package config

import (
	"errors"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	if pathIsValid := fileExists(path); !pathIsValid {
		return nil, errors.New("config path is not valid")
	}

	if err := readFile(path, cfg); err != nil {
		return nil, err
	}

	if err := readEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func readFile(path string, cfg *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func readEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
