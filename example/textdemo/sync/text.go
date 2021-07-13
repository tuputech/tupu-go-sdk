package main

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	textSync "github.com/tuputech/tupu-go-sdk/recognition/text/textsync"
)

func main() {

	// step1. get your secretID
	secretID := "your secretID"
	privateKeyPath := "rsa_private_key.pem"

	// step2. create text handler
	textHandler, err := textSync.NewTextHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// step3. create TextSync object
	textItem := textSync.TextAsyncItem{
		Content: "毛泽东是我们的国家主席",
	}
	texts := make([]textSync.TextAsyncItem, 1)
	texts[0] = textItem

	// start recognition and get result
	result, statusCode, err := textHandler.Perform(secretID, texts)
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

	// Get the value corresponding to the key in json
	code, e = rlt.Get("code").String()
	message, e = rlt.Get("message").String()
	timestamp, e = rlt.Get("timestamp").Int64()
	timestamp = int64(float64(timestamp) / 1000)

	fmt.Printf("- Code: %v %v\n- Time: %v\n", code, message, time.Unix(timestamp, 0))
	for k, v := range task {
		fmt.Printf("- Result: [%v]\n\t%v\n", k, v)
	}
	fmt.Println("- result:", result)
	fmt.Println("----------------------")
}
