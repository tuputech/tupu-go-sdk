package recognition

import "encoding/json"

// Result is a wrapper for service result parsed from response
type Result struct {
	Timestamp int64
	Nonce     string
	Code      int
	Message   string
	Tasks     map[string]interface{}
}

// ParseResult is a helper to parse json string and create a Result struct
func ParseResult(s string) *Result {
	var data map[string]interface{}
	if json.Unmarshal([]byte(s), &data) != nil {
		return nil
	}

	r := Result{}
	r.Tasks = make(map[string]interface{})
	for key, val := range data {
		switch key {
		case "timestamp":
			r.Timestamp = int64(val.(float64))
		case "nonce":
			r.Nonce = val.(string)
		case "code":
			r.Code = int(val.(float64))
		case "message":
			r.Message = val.(string)
		default:
			r.Tasks[key] = val.(map[string]interface{})
		}
	}
	return &r
}
