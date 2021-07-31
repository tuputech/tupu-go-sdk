package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	VDASHdler "github.com/tuputech/tupu-go-sdk/recognition/video/videoasync"
)

func main() {

	// step1. get your secretID
	var (
		// your secretId
		secretID string = "your sid"
		// your rsa_private_key local path
		privateKeyPath string = "your rsa private key"
		// your video url
		videoUrl string = "your need to recognition video url"
		// your receive recognition result server url
		callbackUrl string = "your server url"
		// you will get a videoId after request tupu-video-api
		videoId string
		// json string
		result string
		// video handler
		vdasHdler  *VDASHdler.AsyncHandler
		err        error
		statusCode int
		rlt        *simplejson.Json
	)

	// step2. create speech handler
	vdasHdler, err = VDASHdler.NewVideoAsyncHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}
	// step3. create LongSpeech object
	callbackRules := make(map[string][]VDASHdler.TaskCallbackRule)
	callbackRules["54bcfc6c329af61034f7c2fc"] = []VDASHdler.TaskCallbackRule{
		{
			Label:  1,
			Review: true,
		},
	}

	// start recognition and get result
	// WithXXXX function is optional for api request params
	// e.g. simple to use
	// result, statusCode, err = vdasHdler.Perform(secretID, videoUrl, callbackUrl)
	result, statusCode, err = vdasHdler.Perform(secretID, videoUrl, callbackUrl, VDASHdler.WithCallbackRules(callbackRules))
	printResult(result, statusCode, err)

	// (optional) query recognition result ()
	rlt, err = simplejson.NewJson([]byte(result))

	// get recogniton videoId from Perform func
	videoId, err = rlt.Get("videoId").String()
	if err != nil {
		fmt.Println("start recognition failed, result:", result)
		return
	}

	// (optional) query videoId recognition result
	result, statusCode, err = vdasHdler.QueryRecognitionResult(secretID, videoId)
	printResult(result, statusCode, err)

	// (optional) stop videoId recognition task ()
	result, statusCode, err = vdasHdler.CloseRecognitionTask(secretID, videoId)
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
	task, e = rlt.Map()
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
