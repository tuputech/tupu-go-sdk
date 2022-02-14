package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	SPCHAS "github.com/tuputech/tupu-go-sdk/recognition/speech/speechasync"
)

func main() {

	var (
		// step1. get your secretID
		secretID string = "your secretID"
		// your rsa_private_key local path
		privateKeyPath string = "rsa_private_key.pem"
		// your receive recognition result server url
		callbackUrl string = "your server url"
		// your need to recogniton speech url
		speechUrl string = "your speech url"
		// empty string will using default server url
	)

	// step2. create speech handler
	speechHandler, err := SPCHAS.NewSpeechHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// step3. recognition
	// start recognition and get result
	// WithXXXX function is optional for api request params
	// e.g. simple to use
	// result, statusCode, err := speechHandler.Perform(secretID, speechUrl, callbackUrl)
	result, statusCode, err := speechHandler.Perform(secretID, speechUrl, SPCHAS.WithCallbackURL(callbackUrl), SPCHAS.WithCallbackRule(SPCHAS.CallbackRuleALL), SPCHAS.WithUserId("testId"))
	printResult(result, statusCode, err)
}

func printResult(result string, statusCode int, err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}

	fmt.Println("-------- v1.0 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	// Example of parsing json string using simplejson
	var (
		rlt, e        = simplejson.NewJson([]byte(result))
		task          map[string]interface{}
		code, message string
		timestamp     int64
	)
	if e != nil {
		fmt.Println("[ERROR] params error")
		return
	}
	// fmt.Println(result)

	// Get the value corresponding to the key in json
	code, e = rlt.Get("code").String()
	message, e = rlt.Get("message").String()
	timestamp, e = rlt.Get("timestamp").Int64()
	timestamp = int64(float64(timestamp) / 1000)
	task, e = rlt.Get("result").Map()
	if e != nil {
		fmt.Println("decode error")
		return
	}

	fmt.Printf("- Code: %v %v\n- Time: %v\n", code, message, time.Unix(timestamp, 0))
	for k, v := range task {
		fmt.Printf("- Result: [%v]\n\t%v\n", k, v)
	}
	fmt.Println("----------------------")
}
