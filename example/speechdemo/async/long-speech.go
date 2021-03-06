package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	longSpch "github.com/tuputech/tupu-go-sdk/recognition/speech/longasync"
)

func main() {

	// step1. get your secretID
	secretID := "your secretID"
	privateKeyPath := "rsa_private_key.pem"
	// empty string will using default server url
	serverURL := ""

	// step2. create speech handler
	speechHandler, err := longSpch.NewSpeechHandler(privateKeyPath, serverURL)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}
	// step3. create LongSpeech object
	longSpeech := &longSpch.LongSpeech{
		FileRemoteURL: "your speech url",
		CallbackURL:   "your callback url",
	}

	// start recognition and get result
	result, statusCode, err := speechHandler.Perform(secretID, longSpeech, 0)
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
	fmt.Println("----------------------\n")
}
