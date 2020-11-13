package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/tuputech/tupu-go-sdk/recognition/speech/speechsync"
)

func main() {

	// step1. get your secretID
	secretID := "your secretID"
	privateKeyPath := "your rsa_private_key path"

	// step2. create speech sync handler
	syncHandler, err := speechsync.NewSyncHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// step3. You can choose the request timeout limited, if not, use the default 30s
	// syncHandler.SetTimeout(2)

	// step4. get recognition result

	// test demo1
	testSpeechAPIWithURL(secretID, syncHandler)

	// test demo2
	testSpeechAPIWithPath(secretID, syncHandler)

	// test demo3
	testSpeechAPIWithBinary(secretID, syncHandler)
}

func testSpeechAPIWithBinary(secretID string, speechHandler *speechsync.SyncHandler) {
	//Using local file or binary data
	filePath := "your speech file path"
	fileBytes, e2 := ioutil.ReadFile(filePath)
	if e2 != nil {
		fmt.Printf("Could not load voice: %v", e2)
		return
	}
	// key is your fileName, value is your speech binary data, Extension only supports amr, mp3, wmv, wav, flv
	speechSlice := map[string][]byte{
		"test.amr": fileBytes,
	}

	printResult(speechHandler.PerformWithBinary(secretID, speechSlice))
}

func testSpeechAPIWithPath(secretID string, speechHandler *speechsync.SyncHandler) {
	// step1. get speech file path
	speechPaths := []string{
		"your speech file path",
	}

	// step2. get result of speech recognition API
	printResult(speechHandler.PerformWithPath(secretID, speechPaths))
}

func testSpeechAPIWithURL(secretID string, speechHandler *speechsync.SyncHandler) {
	// step1. get speech file url
	speechURLs := []string{
		"your speech url",
	}
	result, statusCode, err := speechHandler.PerformWithURL(secretID, speechURLs)
	printResult(result, statusCode, err)
}

func printResult(result string, statusCode int, err error) {
	fmt.Printf("result: %s\n", result)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}
	fmt.Println("-------- v1.0 --------")
	// fmt.Println(result)
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

	// Get the value corresponding to the key in json
	code, e = rlt.Get("code").String()
	message, e = rlt.Get("message").String()
	timestamp, e = rlt.Get("timestamp").Int64()
	timestamp = int64(float64(timestamp) / 1000)
	// pase vulgar speech
	task, e = rlt.Get("5c8213b9bc807806aab0a574").Map()
	if e != nil {
		fmt.Println("decode error")
		return
	}

	fmt.Printf("- Code: %v %v\n- Time: %v\n", code, message, time.Unix(timestamp, 0))
	for k, v := range task {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------\n")
}
