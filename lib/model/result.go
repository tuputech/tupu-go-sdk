package model

import (
	"encoding/json"
	"fmt"
)

// Result is a wrapper for service result parsed from response
type Result struct {
	Code      int
	Timestamp int64
	Nonce     string
	Message   string
	Tasks     map[string]interface{}
	Others    map[string]interface{}
}

func parseTaskVal(val interface{}) (v map[string]interface{}, e error) {
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println("Panic info is: ", err)
			e = fmt.Errorf("%v", e)
		}
	}()
	v = val.(map[string]interface{})
	return
}

// ParseResult is a helper to parse json string and create a Result struct
func ParseResult(s string) *Result {

	if len(s) == 0 {
		return nil
	}

	var data map[string]interface{}
	if json.Unmarshal([]byte(s), &data) != nil {
		return nil
	}

	r := Result{}
	r.Tasks = make(map[string]interface{})
	r.Others = make(map[string]interface{})
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
			v, e := parseTaskVal(val)
			if e == nil {
				r.Tasks[key] = v //val.(map[string]interface{})
			} else {
				r.Others[key] = val
			}
		}
	}
	return &r
}
