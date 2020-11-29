package dynamic_test

import (
	"github.com/thmhoag/arkdrater/pkg/arkdrater/dynamic"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"github.com/jinzhu/copier"
)

func Test_RateConverter_GetConvertedRates(t *testing.T) {
	type fields struct {
		http    dynamic.HttpClient
		baseCfg *dynamic.Config
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult *dynamic.Config
		wantErr    bool
	}{
		{
			name: "A single rate is multiplied correctly",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString("TamingSpeedMultiplier=2.0"),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple rates are multiplied correctly",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString(`TamingSpeedMultiplier=2.0
XPMultiplier=2.0
MatingIntervalMultiplier=0.5
EggHatchSpeedMultiplier=2.0
HexagonRewardMultiplier=1.0`),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
					"XPMultiplier": 6,
					"MatingIntervalMultiplier": 0.25,
					"EggHatchSpeedMultiplier": 6,
					"HexagonRewardMultiplier": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Rates not in original config are ignored",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString(`TamingSpeedMultiplier=2.0
HarvestAmountMultiplier=2.0
XPMultiplier=2.0
MatingIntervalMultiplier=0.5
BabyMatureSpeedMultiplier=2.0
EggHatchSpeedMultiplier=2.0
HexagonRewardMultiplier=1.0`),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
					"XPMultiplier": 6,
					"MatingIntervalMultiplier": 0.25,
					"EggHatchSpeedMultiplier": 6,
					"HexagonRewardMultiplier": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Rates with no official override stay the same",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString(`TamingSpeedMultiplier=2.0
BabyMatureSpeedMultiplier=2.0
EggHatchSpeedMultiplier=2.0
HexagonRewardMultiplier=1.0`),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
					"XPMultiplier": 3,
					"MatingIntervalMultiplier": 0.5,
					"EggHatchSpeedMultiplier": 6,
					"HexagonRewardMultiplier": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Non 200 response results error",
			fields: fields{
				http: &HttpClientMock{
					MockGetFunc: func(url string) (resp *http.Response, err error) {
						resp = &http.Response{
							StatusCode: http.StatusBadGateway,
							Status: http.StatusText(http.StatusBadGateway),
							Body: CreateReadCloserFromString(""),
						}

						return
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: nil,
			wantErr: true,
		},
		{
			name: "CRLF line endings are supported",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString("TamingSpeedMultiplier=2.0\r\nXPMultiplier=2.0\r\nMatingIntervalMultiplier=0.5\r\nEggHatchSpeedMultiplier=2.0\r\nHexagonRewardMultiplier=1.0"),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
					"XPMultiplier": 6,
					"MatingIntervalMultiplier": 0.25,
					"EggHatchSpeedMultiplier": 6,
					"HexagonRewardMultiplier": 1,
				},
			},
			wantErr: false,
		},
		{
			name: "LF line endings are supported",
			fields: fields{
				http: &HttpClientMock{
					MockResponse: &http.Response{
						StatusCode: http.StatusOK,
						Body: CreateReadCloserFromString("TamingSpeedMultiplier=2.0\nXPMultiplier=2.0\nMatingIntervalMultiplier=0.5\nEggHatchSpeedMultiplier=2.0\nHexagonRewardMultiplier=1.0"),
					},
				},
				baseCfg: &dynamic.Config{
					Multipliers: map[string]float32{
						"TamingSpeedMultiplier": 5,
						"XPMultiplier": 3,
						"MatingIntervalMultiplier": 0.5,
						"EggHatchSpeedMultiplier": 3,
						"HexagonRewardMultiplier": 1,
					},
				},
			},
			wantResult: &dynamic.Config{
				Multipliers: map[string]float32{
					"TamingSpeedMultiplier": 10,
					"XPMultiplier": 6,
					"MatingIntervalMultiplier": 0.25,
					"EggHatchSpeedMultiplier": 6,
					"HexagonRewardMultiplier": 1,
				},
			},
			wantErr: false,
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyOfOrigCfg := &dynamic.Config{}
			copier.Copy(copyOfOrigCfg, tt.fields.baseCfg)
			r := dynamic.NewRateConverter(tt.fields.http)
			gotResult, err := r.GetConvertedRates(tt.fields.baseCfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConvertedRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("GetConvertedRates() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(tt.fields.baseCfg, copyOfOrigCfg) {
				t.Errorf("GetConvertedRates() should not modify the original config")
			}
		})
	}
}

type HttpClientMock struct {
	// MockGetFunc is used to mock the http.Get function
	MockGetFunc func(url string) (resp *http.Response, err error)
	// MockResponse will be returned if the MockGetFunc is not set
	MockResponse *http.Response
	// MockErr will be returned if the MockGetFunc is not set
	MockErr error
}

func (mock HttpClientMock) Get(url string) (resp *http.Response, err error) {
	if mock.MockGetFunc!= nil {
		return mock.MockGetFunc(url)
	}

	return mock.MockResponse, mock.MockErr
}

func CreateReadCloserFromString(str string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(str))
}