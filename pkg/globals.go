package h0neytr4p

import (
	"os"
	"sync"
)

type Trap struct {
	Basicinfo struct {
		Name            string `json:"Name"`
		Port            string `json:"Port"`
		Protocol        string `json:"Protocol"`
		Mitreattacktags string `json:"MitreAttackTags"`
		RiskRating      string `json:"RiskRating"`
		References      string `json:"References"`
		Description     string `json:"Description"`
	} `json:"BasicInfo"`
	Behaviour []struct {
		Request struct {
			URL     string                 `json:"Url"`
			Method  string                 `json:"Method"`
			Headers map[string]interface{} `json:"Headers"`
			Params  map[string]interface{} `json:"Params"`
		} `json:"Request"`
		Response struct {
			Statuscode int                    `json:"StatusCode"`
			Body       string                 `json:"Body"`
			Headers    map[string]interface{} `json:"Headers"`
			Type       string                 `json:"Type"`
		} `json:"Response"`
		Trap string `json:"trap,omitempty"`
	} `json:"Behaviour"`
}

const (
	MaxMultipartSize = 101 * 1024 // 101KB
	MaxJSONFormSize  = 11 * 1024  // 11KB
)

var (
	logFile       *os.File
	logFileMutex  sync.Mutex
	payloadFolder string
	Verbose       string
)
