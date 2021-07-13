package main

import (
	"fmt"
	"io/ioutil"

	//rcn "recognition"
	"time"

	rcn "github.com/tuputech/tupu-go-sdk/recognition"
)

func main() {
	secretID := "Your SecretID"
	handler, e := rcn.NewHandler("rsa_private_key.pem")
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}

	//Optional Step: using http-client created by your own
	// tr := &http.Transport{
	// 	MaxIdleConns:       10,
	// 	IdleConnTimeout:    30 * time.Second,
	// 	DisableCompression: true,
	// }
	// handler.Client = &http.Client{Transport: tr}

	images1 := []string{"your image url"}

	// just for images
	printResult(handler.PerformWithURL(secretID, images1))

	// Number of tags less than number of images, the rest images will use the last tag
	tags := []string{"image tag"}
	printResult(handler.PerformWithURL(secretID, images1, handler.WithTags(tags)))

	// run by appoint task
	tasks := []string{"54bcfc6c329af61034f7c2fc"}
	printResult(handler.PerformWithURL(secretID, images1, handler.WithTasks(tasks)))

	// with tag and run by appoint task
	printResult(handler.PerformWithURL(secretID, images1, handler.WithTasks(tasks), handler.WithTags(tags)))

	// Using local file path
	filepaths := []string{"your filepath"}
	printResult(handler.PerformWithPath(secretID, filepaths, handler.WithTags(tags)))

	// Using local file or binary data
	fileBytes, e2 := ioutil.ReadFile("your speech filePath")
	if e2 != nil {
		fmt.Printf("Could not load image: %v", e2)
		return
	}
	imgBinary := rcn.NewBinaryImage(fileBytes, "1.jpg")
	defer imgBinary.ClearBuffer()
	images2 := []*rcn.Image{rcn.NewLocalImage("your speech filePath"), imgBinary}
	printResult(handler.Perform(secretID, images2, []string{"Local Image", "Using Buffer"}, nil))
}

func printResult(result string, statusCode int, e error) {
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	fmt.Println("-------- v1.6 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := rcn.ParseResult(result)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(int64(float64(r.Timestamp)/1000), 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------")
}
