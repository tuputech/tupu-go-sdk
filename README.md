# TUPU GO SDK

Golang SDK for TUPU visual recognition service (v1.4)
######  
<https://www.tuputech.com>

## Changelogs
#### v1.4
- removed log.Fatal

#### v1.3
- fixed bug in parsing result

#### v1.2
- add shortcut methods for URL or path

#### v1.1
- 1st ready version

## Installation
go get github.com/tuputech/tupu-go-sdk/recognition

## Example

```
package main

import (
	"fmt"
	"io/ioutil"
	//rcn "recognition"
	"time"

	rcn "github.com/tuputech/tupu-go-sdk/recognition"
)

func main() {
	secretID := "your-secret-id"
	handler, e := rcn.NewHandler("rsa_private_key.pem")
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	//Optional Step: set identity of sub-user if necessary
	//handler.UID = "bucket-of-jackbauer"

	url1 := "http://www.yourdomain.com/img/1.jpg"
	url2 := "http://www.yourdomain.com/img/2.jpg"
	images1 := []string{url1, url2}
	//No tag for images
	printResult(handler.PerformWithURL(secretID, images1, nil))
	//Number of tags less than number of images, the rest images will use the last tag
	printResult(handler.PerformWithURL(secretID, images1, []string{"Remote Image"}))

	//Using local file or binary data
	fileBytes, e2 := ioutil.ReadFile("img/1.jpg")
	if e2 != nil {
		fmt.Printf("Could not load image: %v", e2)
		return
	}
	imgBinary := rcn.NewBinaryImage(fileBytes, "1.jpg")
	defer imgBinary.ClearBuffer()
	images2 := []*rcn.Image{rcn.NewLocalImage("img/2.jpg"), imgBinary}
	printResult(handler.Perform(secretID, images2, []string{"Local Image", "Using Buffer"}))
}

func printResult(result string, statusCode int, e error) {
	if e != nil {
		fmt.Printf("Failed: %v\n", e)
		return
	}
	fmt.Println("-------- v1.2 --------")
	fmt.Printf("Status-Code: %v\n-----\n", statusCode)

	r := rcn.ParseResult(result)
	fmt.Printf("- Code: %v %v\n- Time: %v\n", r.Code, r.Message, time.Unix(r.Timestamp, 0))
	for k, v := range r.Tasks {
		fmt.Printf("- Task: [%v]\n%v\n", k, v)
	}
	fmt.Println("----------------------")
}

```

----------------------

### func PerformWithURL
func (h *Handler) PerformWithURL(secretID string, imageURLs []string, tags []string) (result string, statusCode int, e error)

- **secretId**: secret-id for recognition tasks
- **imageURLs**: array of image URLs
- **tags**: array of tags for images (optional)

----------------------

### func PerformWithPath
func (h *Handler) PerformWithPath(secretID string, imagePaths []string, tags []string) (result string, statusCode int, e error)

- **secretId**: secret-id for recognition tasks
- **imagePaths**: array of image paths
- **tags**: array of tags for images (optional)

----------------------

### func Perform
func (h *Handler) Perform(secretID string, images []*Image, tags []string) (result string, statusCode int, e error)

- **secretId**: secret-id for recognition tasks
- **images**: array of Image struct, but don't mix use of URL and path/binary in one call
- **tags**: array of tags for images (optional)

## License

[MIT](http://www.opensource.org/licenses/mit-license.php)
