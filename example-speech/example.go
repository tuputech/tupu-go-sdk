package main

import (
	"fmt"
	"io/ioutil"
	"time"

	spch "github.com/tuputech/tupu-go-sdk/recognition-api/speech/shortsync"
)

func main() {

	// step1. get your secretID
	secretID := "5f042c1f1bac63001e897f27"
	// step2. create speech handler
	speechHandler, err := spch.NewSpeechHandlerWithURL("/Users/mac/hcz/go_project/tupu_rsa_key/rsa_private_key.pem",
		"http://172.26.2.63:8991/v3/recognition")
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// step3. get recognition result
	// test demo1
	testSpeechAPIWithURL(secretID, speechHandler)

	// test demo2
	//testSpeechAPIWithPath(secretID, speechHandler)

	// test demo3
	//testSpeechAPIWithBinary(secretID, speechHandler)
}

func testSpeechAPIWithBinary(secretID string, speechHandler *spch.SpeechHandler) {
	//Using local file or binary data
	fileBytes, e2 := ioutil.ReadFile("/Users/mac/Music/vulgar.wmv")
	if e2 != nil {
		fmt.Printf("Could not load voice: %v", e2)
		return
	}
	speechBinary := spch.NewBinarySpeech(fileBytes, "1.wmv")
	defer speechBinary.ClearBuffer()
	speechSlice := []*spch.Speech{spch.NewLocalSpeech("/Users/mac/Music/vulgar.wmv"), speechBinary}

	printResult(speechHandler.Perform(secretID, speechSlice))
}

func testSpeechAPIWithPath(secretID string, speechHandler *spch.SpeechHandler) {
	// step1. get speech file path
	speechPaths := []string{
		"/Users/mac/Music/vulgar.wmv",
	}

	// step2. get result of speech recognition API
	printResult(speechHandler.PerformWithPath(secretID, speechPaths))
}

func testSpeechAPIWithURL(secretID string, speechHandler *spch.SpeechHandler) {
	// step1. get speech file url
	speechURLs := []string{
		"https://www.tuputech.com/original/world/data-c40/yrw/api_test_data/vulgar.wmv",
	}
	printResult(speechHandler.PerformWithURL(secretID, speechURLs))
}

func printResult(result string, statusCode int, err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}
	fmt.Println("-------- v1.0 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := spch.ParseResult(result)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(r.Timestamp, 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------\n")
}
