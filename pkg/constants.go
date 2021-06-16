package h0neytr4p

import "os"

type LogEntry struct {
	SourceIP  string
	UserAgent string
	Timestamp string
	Path      string
	Trapped   string
}
type Attacker struct {
	SourceIP   string `json:"SourceIP"`
	UserAgent  string `json:"UserAgent"`
	TrappedFor string `json:"TrappedFor"`
	RiskRating string `json:"RiskRating"`
	References string `json:"References"`
	Timestamp  string `json:"Timestamp"`
}

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

var filenameGlobal *os.File
var logFileGlobal *os.File
var Verbose string
