package config

import "github.com/thmhoag/arkdrater/pkg/arkdrater/dynamic"

type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	DynamicConfig dynamic.Config `yaml:"dynamicConfig"`
}
