package handler

import (
	"github.com/thmhoag/arkdrater/pkg/arkdrater/config"
	"github.com/thmhoag/arkdrater/pkg/arkdrater/dynamic"
	"io"
	"log"
	"net/http"
)

type DefaultHandler struct {
	RateConverter dynamic.RateConverter
	Config config.Config
}

func (h *DefaultHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	convertedRatesCfg, err:= h.RateConverter.GetConvertedRates()
	if err != nil {
		log.Printf("error converting rates using official multipliers: %v\n", err)
		writeCfg(w, &h.Config.DynamicConfig)
	}

	writeCfg(w, convertedRatesCfg)
}

func writeCfg(w io.Writer, cfg *dynamic.Config) {
	if err := cfg.WriteIniStr(w); err != nil {
		log.Printf("error writing config values: %v\n", err)
	}
}