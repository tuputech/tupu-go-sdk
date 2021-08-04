package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/bitly/go-simplejson"
	VDSHdler "github.com/tuputech/tupu-go-sdk/recognition/video/videosync"
)

func main() {

	// step1. get your secretID
	var (
		// your secretId
		secretID string
		// your rsa_private_key local path
		privateKeyPath string
		//
		// video handler
		vdsHdler *VDSHdler.SyncHandler
		err      error
	)

	// step2. create speech handler
	vdsHdler, err = VDSHdler.NewSyncHandler(privateKeyPath)
	if err != nil {
		fmt.Println("-------- ERROR ----------")
		return
	}

	// video recognition with resource url
	testVideoSyncAPIWithURL(secretID, vdsHdler)

	// video recognition with resource local path
	testVideoSyncAPIWithPath(secretID, vdsHdler)

	// video recognition with resource binary data
	testVideoSyncAPIWithBinary(secretID, vdsHdler)

}

func testVideoSyncAPIWithURL(secretID string, vdsHdler *VDSHdler.SyncHandler) {
	var (
		result     string
		statusCode int
		err        error
		videoUrls  = []string{
			"recognition video url",
		}
	)

	// WithXXXX function is optional for api request params
	// e.g. simple to use
	// result, statusCode, err = vdsHdler.PerformWithURL(secretID, videoUrls)
	result, statusCode, err = vdsHdler.PerformWithURL(secretID, videoUrls, VDSHdler.WithInterval(5), VDSHdler.WithMaxFrames(200), VDSHdler.WithTag("my recognition recognition"))
	printResult(result, statusCode, err)
}

func testVideoSyncAPIWithPath(secretID string, vdsHdler *VDSHdler.SyncHandler) {
	var (
		result     string
		statusCode int
		err        error
		videoPaths = []string{
			"recognition local video path",
		}
	)

	// WithXXXX function is optional params for api request params
	result, statusCode, err = vdsHdler.PerformWithPath(secretID, videoPaths, VDSHdler.WithInterval(5))
	printResult(result, statusCode, err)
}

func testVideoSyncAPIWithBinary(secretID string, vdsHdler *VDSHdler.SyncHandler) {
	//Using local file or binary data
	var (
		filePath       = "your speech filePath"
		fileBytes, err = ioutil.ReadFile(filePath)
		result         string
		statusCode     int
	)
	if err != nil {
		fmt.Printf("Could not load voice: %v", err)
		return
	}
	// key is your fileName, value is the speech binary data
	videoBinarys := map[string][]byte{
		filepath.Base(filePath): fileBytes,
	}

	// WithXXXX function is optional params for api request params
	result, statusCode, err = vdsHdler.PerformWithBinary(secretID, videoBinarys, VDSHdler.WithInterval(10))
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
