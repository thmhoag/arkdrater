package dynamic

type Config struct {
	Multipliers                  map[string]float32 `yaml:"multipliers"`
	TributeItemExpirationSeconds int64              `yaml:"TributeItemExpirationSeconds" envconfig:"TributeItemExpirationSeconds"`
	TributeDinoExpirationSeconds int64              `yaml:"TributeDinoExpirationSeconds" envconfig:"TributeDinoExpirationSeconds"`
	EnableFullDump               bool               `yaml:"EnableFullDump" envconfig:"EnableFullDump"`
	GUseServerNetSpeedCheck      bool               `yaml:"GUseServerNetSpeedCheck" envconfig:"GUseServerNetSpeedCheck"`
	bUseAlarmNotifications       bool               `yaml:"bUseAlarmNotifications" envconfig:"bUseAlarmNotifications"`
}