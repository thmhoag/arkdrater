package dynamic

import (
	"io"
	"log"
	"text/template"
)

type Config struct {
	Multipliers                  map[string]float32 `yaml:"multipliers"`
	TributeItemExpirationSeconds int64              `yaml:"TributeItemExpirationSeconds" envconfig:"TributeItemExpirationSeconds"`
	TributeDinoExpirationSeconds int64              `yaml:"TributeDinoExpirationSeconds" envconfig:"TributeDinoExpirationSeconds"`
	EnableFullDump               bool               `yaml:"EnableFullDump" envconfig:"EnableFullDump"`
	GUseServerNetSpeedCheck      bool               `yaml:"GUseServerNetSpeedCheck" envconfig:"GUseServerNetSpeedCheck"`
	bUseAlarmNotifications       bool               `yaml:"bUseAlarmNotifications" envconfig:"bUseAlarmNotifications"`
}

func (c *Config) WriteIniStr(w io.Writer) error {
	template, err := template.New("handlerResponse").Parse(iniTemplate)
	if err != nil {
		log.Printf("error while marshling object. %s \n", err.Error())
		return err
	}

	if err := template.Execute(w, c); err != nil {
		log.Printf("template.Execute error %s \n", err.Error())
		return err
	}

	return nil
}

const iniTemplate = `
{{- range $multiplier, $value := .Multipliers }}
{{ $multiplier }}={{ printf "%.1f" $value }}
{{- end }}
`