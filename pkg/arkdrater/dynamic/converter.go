package dynamic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	// OfficialDynamicConfigUrl is the url to Studio Wildcard's official dynamic config
	OfficialDynamicConfigUrl string = "https://cdn2.arkdedicated.com/asa/dynamicconfig.ini"
)

// HttpClient is an interface for the official http.Client
type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type rateConv struct {
	http HttpClient
}

// RateConverter can convert base config rates to special event rates
type RateConverter interface {
	// GetConvertedRates multiplies the base rates by the current official event rates
	GetConvertedRates(baseCfg *Config) (result *Config, err error)
}

// NewRateConverter returns a new RateConverter with the specified base config
func NewRateConverter(http HttpClient) RateConverter {
	if http == nil {
		panic("http cannot be nil")
	}

	return &rateConv{http: http}
}

func (r rateConv) GetConvertedRates(baseCfg *Config) (result *Config, err error) {
	result = baseCfg.GetCopy()
	resp, err := r.http.Get(OfficialDynamicConfigUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response for official rates: %v\n%v\n", resp.StatusCode, resp.Status)
	}

	offDynaConfBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	offDynaConfMap := make(map[string]float32)
	offDynaConfLines := strings.Split(string(offDynaConfBody), "\n")
	for _, line := range offDynaConfLines {
		// make sure we remove artifacts from Windows line endings
		cleansedLine := strings.ReplaceAll(line, "\r", "")
		splitLine := strings.Split(cleansedLine, "=")
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

	for name, value := range result.Multipliers {
		if officialValue, ok := offDynaConfMap[name]; ok {
			newValue := value * officialValue
			result.Multipliers[name] = newValue
		}
	}

	return result, nil
}

