package longasync

import (
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

// SpeechResult is a wrapper for service result parsed from response
type SpeechResult struct {
	*tupumodel.Result
}

// ParseResult is a helper to parse json string and create a Result struct
func ParseResult(result string) *SpeechResult {
	if len(result) == 0 {
		return nil
	}
	speechRlt := new(SpeechResult)
	speechRlt.Result = tupumodel.ParseResult(result)
	return speechRlt
}
