package shortsync

import (
	generalrcn "github.com/tuputech/tupu-go-sdk/recognition-api/general"
)

// SpeechResult is a wrapper for service result parsed from response
type SpeechResult struct {
	generalrcn.Result
}

// ParseResult is a helper to parse json string and create a Result struct
func ParseResult(result string) *SpeechResult {
	speechRlt := new(SpeechResult)
	speechRlt.Result = *generalrcn.ParseResult(result)
	return speechRlt
}
