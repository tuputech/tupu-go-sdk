package main

import (
	"fmt"
	"time"

	longSpch "github.com/tuputech/tupu-go-sdk/recognition-api/speech/longasync"
)

func main() {

	// step1. get your secretID
	secretID := "5f48703b9b55e9001e694707"

	// step2. create speech handler
	speechHandler, err := longSpch.NewSpeechHandlerWithURL(
		"/Users/mac/hcz/go_project/tupu_rsa_key/rsa_private_key.pem",
		"http://172.26.2.63:51021/v3/recognition/speech/recording/async/",
	)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}
	// step3. create LongSpeech object
	longSpeech := &longSpch.LongSpeech{
		URL:         "https://www.tuputech.com/original/world/data-c40/yrw/api_test_data/vulgar.wmv",
		CallbackURL: "http:172.26.2.63:8991",
	}

	// start recognition and get result
	printResult(speechHandler.Perform(secretID, longSpeech))
}

func printResult(result string, statusCode int, err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}
	fmt.Println("-------- v1.0 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := longSpch.ParseResult(result)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(r.Timestamp, 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------\n")
}
