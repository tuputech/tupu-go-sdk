package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/tuputech/tupu-go-sdk/recognition/speech/speechasync"
)

func main() {

	// step1. get your secretID
	secretID := "your secretID"
	privateKeyPath := "your rsa_private_key path"

	// step2. create speech handler
	speechHandler, err := speechasync.NewSpeechHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}
	// step3. You can choose the request timeout limited, if not, use the default 30s
	// speechHandler.SetTimeout(2)

	// step4. create LongSpeech object
	speechAsync := &speechasync.SpeechAsync{
		FileRemoteURL: "your speech url",
		CallbackURL:   "your callback url",
	}

	// start recognition and get result
	result, statusCode, err := speechHandler.Perform(secretID, speechAsync)
	printResult(result, statusCode, err)
}

func printResult(result string, statusCode int, err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}

	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	// Example of parsing json string using simplejson
	var (
		rlt, e          = simplejson.NewJson([]byte(result))
		task            map[string]interface{}
		message         string
		timestamp, code int64
	)
	if e != nil {
		fmt.Println("[ERROR] params error")
		return
	}
	// fmt.Println(result)

	// Get the value corresponding to the key in json
	code, e = rlt.Get("code").Int64()
	message, e = rlt.Get("message").String()
	timestamp, e = rlt.Get("timestamp").Int64()
	timestamp = int64(float64(timestamp) / 1000)
	task, _ = rlt.Get("result").Map()
	if (e != nil) || (code != 0) {
		fmt.Println("decode error, code: ", code)
		fmt.Println("message: ", message)
		return
	}

	fmt.Printf("- Code: %v %v\n- Time: %v\n", code, message, time.Unix(timestamp, 0))
	for k, v := range task {
		fmt.Printf("- Result: [%v]\n\t%v\n", k, v)
	}
	fmt.Println("----------------------\n")
}
