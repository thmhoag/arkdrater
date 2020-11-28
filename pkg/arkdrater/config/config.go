package config

type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	DynamicConfig struct {
		Multipliers                  map[string]float32 `yaml:"multipliers"`
		TributeItemExpirationSeconds int64              `yaml:"TributeItemExpirationSeconds" envconfig:"arkdrater_TributeItemExpirationSeconds"`
		TributeDinoExpirationSeconds int64              `yaml:"TributeDinoExpirationSeconds" envconfig:"arkdrater_TributeDinoExpirationSeconds"`
		EnableFullDump               bool               `yaml:"EnableFullDump" envconfig:"arkdrater_EnableFullDump"`
		GUseServerNetSpeedCheck      bool               `yaml:"GUseServerNetSpeedCheck" envconfig:"arkdrater_GUseServerNetSpeedCheck"`
		bUseAlarmNotifications       bool               `yaml:"bUseAlarmNotifications" envconfig:"arkdrater_BbUseAlarmNotifications"`
	} `yaml:"dynamicConfig"`
}
