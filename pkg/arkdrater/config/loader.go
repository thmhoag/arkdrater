package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/thmhoag/arkdrater/pkg/arkdrater/dynamic"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const envMultiPrefix = "arkdratermulti_"

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		DynamicConfig: dynamic.Config{
			Multipliers: map[string]float32{},
		},
	}

	if pathIsValid := fileExists(path); !pathIsValid {
		return nil, errors.New("config path is not valid")
	}

	if err := readFile(path, cfg); err != nil {
		return nil, err
	}

	if err := readEnv(cfg); err != nil {
		return nil, err
	}

	readMultiplierEnv(cfg)

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

func readMultiplierEnv(cfg *Config) {
	for _, element := range os.Environ() {
		env := strings.Split(element, "=")
		if len(env) != 2 {
			// not valid
			continue
		}

		name := strings.TrimSpace(env[0])
		if name == "" {
			// not valid
			continue
		}


		if !strings.HasPrefix(name, envMultiPrefix) {
			continue
		}

		name = strings.Replace(name, envMultiPrefix, "", 1)
		name = strings.TrimSpace(name)
		if name == "" {
			// not valid
			continue
		}

		val, err := strconv.ParseFloat(env[1], 32)
		if err != nil {
			log.Printf("%v%v does not have a valid float value (found: %v)\n", envMultiPrefix, name, env[1])
			continue
		}

		cfg.DynamicConfig.Multipliers[name] = float32(val)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
