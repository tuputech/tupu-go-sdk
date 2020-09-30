package main

import (
	"fmt"
	"io/ioutil"

	//rcn "recognition"
	"time"

	rcn "github.com/tuputech/tupu-go-sdk/recognition"
)

func main() {
	secretID := "5d3982e0d06eeb0facac6c58"
	handler, e := rcn.NewHandler("rsa_private_key.pem")
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	//Optional Step: set identity of sub-user if necessary
	//handler.UID = "bucket-of-jackbauer"

	//Optional Step: using http-client created by your own
	// tr := &http.Transport{
	// 	MaxIdleConns:       10,
	// 	IdleConnTimeout:    30 * time.Second,
	// 	DisableCompression: true,
	// }
	// handler.Client = &http.Client{Transport: tr}

	url1 := "https://www.tuputech.com/original/world/data-c40/yrw/api_test_data/normal.jpg"
	url2 := "http://www.yourdomain.com/img/2.jpg"
	images1 := []string{url1, url2}
	//No tag for images
	printResult(handler.PerformWithURL(secretID, images1, nil))
	//Number of tags less than number of images, the rest images will use the last tag
	// printResult(handler.PerformWithURL(secretID, images1, []string{"Remote Image"}))

	//Using local file or binary data
	fileBytes, e2 := ioutil.ReadFile("/Users/mac/Downloads/tp.jpeg")
	if e2 != nil {
		fmt.Printf("Could not load image: %v", e2)
		return
	}
	imgBinary := rcn.NewBinaryImage(fileBytes, "1.jpg")
	defer imgBinary.ClearBuffer()
	// images2 := []*rcn.Image{rcn.NewLocalImage("/Users/mac/Downloads/tp.jpeg"), imgBinary}
	// printResult(handler.Perform(secretID, images2, []string{"Local Image", "Using Buffer"}))
}

func printResult(result string, statusCode int, e error) {
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	fmt.Println("-------- v1.6 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := rcn.ParseResult(result)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(r.Timestamp, 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------")
}
