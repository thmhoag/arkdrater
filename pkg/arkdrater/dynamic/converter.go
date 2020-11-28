package dynamic

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	// OfficialDynamicConfigUrl is the url to Studio Wildcard's official dynamic config
	OfficialDynamicConfigUrl string = "http://arkdedicated.com/dynamicconfig.ini"
)

// HttpClient is an interface for the official http.Client
type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type rateConv struct {
	http HttpClient
	cfg Config
}

// RateConverter can convert base config rates to special event rates
type RateConverter interface {
	// GetConvertedRates multiplies the base rates by the current official event rates
	GetConvertedRates() (result *Config, err error)
}

// NewRateConverter returns a new RateConverter with the specified base config
func NewRateConverter(http HttpClient, cfg Config) RateConverter {
	if http == nil {
		panic("http cannot be nil")
	}

	return &rateConv{http: http, cfg: cfg}
}

func (r rateConv) GetConvertedRates() (result *Config, err error) {
	newCfg := r.cfg
	result = &newCfg

	resp, err := r.http.Get(OfficialDynamicConfigUrl)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	offDynaConfBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	offDynaConfMap := make(map[string]float32)
	offDynaConfLines := strings.Split(string(offDynaConfBody), "\n")
	for _, line := range offDynaConfLines {
		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			// error?
			continue
		}

		multName := strings.TrimSpace(splitLine[0])
		if multName == "" {
			// invalid, do something?
			continue
		}

		multVal, err := strconv.ParseFloat(splitLine[1], 32)
		if err != nil {
			// do something?
			continue
		}

		offDynaConfMap[multName] = float32(multVal)
	}

	for name, value := range r.cfg.Multipliers {
		if officialValue, ok := offDynaConfMap[name]; ok {
			newValue := value * officialValue
			result.Multipliers[name] = newValue
		}
	}

	return result, nil
}

