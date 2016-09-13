package main

import (
	"fmt"
	"io/ioutil"
	rcn "recognition"
	"time"
)

func main() {
	handler, e := rcn.NewHandler("rsa_private_key.pem")
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	//Set identity of sub-user if necessary
	handler.UID = "bucket-of-jackbauer"

	//Tag is optional
	url := "http://www.yourdomain.com/img/1.jpg"
	imgLink := rcn.NewRemoteImage(url).Tag("Remote Image")
	imgPath := rcn.NewLocalImage("img/2.jpg").Tag("Local Image")
	fileBytes, e2 := ioutil.ReadFile("img/1.jpg")
	if e2 != nil {
		fmt.Printf("Could not load image: %v", e2)
		return
	}
	imgBinary := rcn.NewBinaryImage(fileBytes, "1.jpg")
	imgBinary.Tag("Using Buffer")
	defer imgBinary.ClearBuffer()

	images1 := []*rcn.Image{imgLink}
	images2 := []*rcn.Image{imgPath, imgBinary}

	secretID := "your-secret-id"
	send(handler, secretID, images1)
	send(handler, secretID, images2)
}

func send(h *rcn.Handler, secretID string, images []*rcn.Image) {
	json, statusCode, e := h.Perform(secretID, images)
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	fmt.Println("---------------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := rcn.ParseResult(json)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(r.Timestamp, 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("---------------")
}
