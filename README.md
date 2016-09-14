# TUPU GO SDK

Golang SDK for TUPU visual recognition service (v1.1)
######  
<https://www.tuputech.com>

## Changelogs
#### v1.1
- 1st ready version

## Example

```
package main

import (
	"fmt"
	"io/ioutil"
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

	images1 := []*rcn.Image{rcn.NewRemoteImage(url1), rcn.NewRemoteImage(url2)}
	//No tag for images
	send(handler, secretID, images1, nil)
	//Number of tags less than number of images, the rest images will use the last tag
	send(handler, secretID, images1, []string{"Remote Image"})

	//Using local file or binary data
	fileBytes, e2 := ioutil.ReadFile("img/1.jpg")
	if e2 != nil {
		fmt.Printf("Could not load image: %v", e2)
		return
	}
	imgBinary := rcn.NewBinaryImage(fileBytes, "1.jpg")
	defer imgBinary.ClearBuffer()
	images2 := []*rcn.Image{rcn.NewLocalImage("img/2.jpg"), imgBinary}
	send(handler, secretID, images2, []string{"Local Image", "Using Buffer"})
}

func send(h *rcn.Handler, secretID string, images []*rcn.Image, tags []string) {
	json, statusCode, e := h.Perform(secretID, images, tags)
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
```

### func Perform
func (h *Handler) Perform(secretID string, images []*Image, tags []string) (result string, statusCode int, e error)

- **$secretId**: secret-id for recognition tasks
- **$images**: array of image URLs or paths or file binary (don't mix use of URL and path or binary in one call)
- **$tags**: array of tags for images (optional)

## License

[MIT](http://www.opensource.org/licenses/mit-license.php)
