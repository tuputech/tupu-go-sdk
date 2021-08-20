package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	SPSTRM "github.com/tuputech/tupu-go-sdk/recognition/speech/speechstream"
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
		streamUrl string = "your speech url"
		// empty string will using default server url
		rlt        *simplejson.Json
		requestId  string
		result     string
		statusCode int
		err        error
	)

	// step2. create speech handler
	spstrmHandler, err := SPSTRM.NewSpeechStreamHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// step3. start recognition
	// start recognition and get result
	// WithXXXX function is optional for api request params
	// e.g. simple to use
	// result, statusCode, err := spstrmHandler.Perform(secretID, speechUrl, callbackUrl)
	result, statusCode, err = spstrmHandler.StartStreamRecognition(secretID, streamUrl, callbackUrl,
		SPSTRM.WithCallbackRules(SPSTRM.CallbackAllRecognition),
		SPSTRM.WithTask(SPSTRM.SpeechAnalysisTaskID),
	)

	// step4. parse response body to get requestId
	rlt, err = simplejson.NewJson([]byte(result))
	requestId, err = rlt.Get("result").GetIndex(0).Get("requestId").String()
	if err != nil {
		fmt.Println("start recognition failed, result:", result)
		return
	}
	printResult(result, statusCode, err)

	// step5. close recognition task
	result, statusCode, err = spstrmHandler.CloseRecognitionTask(secretID, requestId)
	fmt.Println("close recognition task result: ", result, "\nstatus code: ", statusCode, "\nerror: ", err)

	// step6. (optional) query recognition status
	result, statusCode, err = spstrmHandler.QueryStatus(secretID, requestId)
	fmt.Println("query recognition status result: ", result, "\nstatus code: ", statusCode, "\nerror: ", err)
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
