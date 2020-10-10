# TUPU GO SDK

Golang SDK for TUPU visual recognition service (v1.6.1)
######  
<https://www.tuputech.com>

## Changelogs
#### v1.7
- add speech method and lib

#### v1.6.1
- fix returning data when it's not JSON data

#### v1.6
- fix to return failure response status code

#### v1.5
- support setting http client

#### v1.4
- removed log.Fatal

#### v1.3
- fixed bug in parsing result

#### v1.2
- add shortcut methods for URL or path

#### v1.1
- 1st ready version

## Installation

​	go get github.com/tuputech/tupu-go-sdk/recognition

## Example

1. [Image recognition interface example](./example/image.go)  
2. [shortSpeech recognition interface example](./example/short-speech.go)  
3. [longSpeech recognition interface example](./example/long-speech.go)  

## Image Recognition API

> import "github.com/tuputech/recognition"

---

### func <font color=#71ABD5>PerformWithURL</font>

```go  
func (h *Handler) PerformWithURL(secretID string, imageURLs []string, tags []string) (result string, statusCode int, e error)
```
**PerformWithURL** return a json `string` and a `int` express response, a `error` identifying success of failure

It is useful for the recognition function using remote image  

   > Params Descrition

   - ***secretId***: secret-id for recognition tasks
   - ***imageURLs***: array of image URLs
   - ***tags***: array of tags for images (optional)

---

### func <font color=#71ABD5>PerformWithPath</font>
```go  
func (h *Handler) PerformWithPath(secretID string, imagePaths []string, tags []string) (result string, statusCode int, e error)
```
**PerformWithPath** return a json `string` and a `int` express response, a `error` identifying success of failure

It is useful for the recognition function using local image  

> **Params Descrition**  
- ***secretId***: secret-id for recognition tasks
- ***imagePaths***: array of image paths
- ***tags***: array of tags for images (optional)

---

### func <font color=#71ABD5>Perform</font>

```go
func (h *Handler) Perform(secretID string, images []*Image, tags []string) (result string, statusCode int, e error)
```

**Perform** return a json `string` and a `int` express response, a `error` identifying success of failure

Construct the data structures we provide to execute reccognition

There are three functions you can use to construct an `Image` object：

1. `func NewRemoteImage(url string) *Image`
2. `func NewLocalImage(path string) *Image `
3. `func NewBinaryImage(buf []byte, filename string) *Image`

> **Params  Descrition**
- ***secretId***: secret-id for recognition tasks
- ***images***: array of Image struct, but don't mix use of URL and path/binary in one call
- ***tags***: array of tags for images (optional)

---

## Speech Recognition API

> Contains `Package shortsync` and `Package longasync`

### short speech API

> import "github.com/tuputech/recognition/speech/shortsync"

---

#### func <font color=#71ABD5>PerforWithBinary</font>

```go
func (spHdler *ShortSpeechHandler) PerformWithBinary(secretID string, binaryData map[string][]byte, timeout int) (result string, statusCode int, err error)
```

**PerformWithBinary** return a json `string` and a `int` express response, a `error` identifying success of failure

Identification with binaries is valid, but binaries need to be built with the Map type for filetype supported only in `amr, mp3, wmv, wav, flv` format

> **Params  Descrition**

- ***secretID***: secret-id for recognition task
- ***binaryData***: map type, key means file name, value means binary data
- ***timeout***: Set request timeout, if value equal 0, will using default timeout(30s)

---

#### func <font color=#71ABD5>PerforWithPath</font>

```go
func (spHdler *ShortSpeechHandler) PerformWithPath(secretID string, speechPaths []string, timeout int) (result string, statusCode int, err error)
```

**PerformWithPath** return a json `string` and a `int` express response, a `error` identifying success of failure

It is useful for the recognition function using local speech file

> **Params  Descrition**

- ***secretID***: secret-id for recognition task
- ***speechPaths***: local short speech paths
- ***timeout***: Set request timeout, if value equal 0, will using default timeout(30s)

-----

#### func <font color=#71ABD5>PerforWithURL</font>

```go
func (spHdler *ShortSpeechHandler) PerformWithURL(secretID string, URLs []string, timeout int) (result string, statusCode int, err error)
```

**PerformWithURL** return a json `string` and a `int` express response, a `error` identifying success of failure

It is useful for the recognition function using remote speech file 

> **Params  Descrition**

- ***secretID***: secret-id for recognition task
- ***URLs***: remote short speech address
- ***timeout***: Set request timeout, if value equal 0, will using default timeout(30s)

---

### long speech API

> import "github.com/tuputech/recognition/speech/longasync"

---

#### func <font color=#71ABD5>Perform</font>

```go
func (spHdler *LongSpeechHandler) Perform(secretID string, longspch *LongSpeech, timeout int) (result string, statusCode int, err error)
```

 ```go
// LongSpeech is a structure that encapsulates long speech messages
type LongSpeech struct {
	// FileRemoteURL represents the address of the long voice, can't be empty
	FileRemoteURL string 
	// CallbackUrl represents the address of the callback result, cant' be empty
	CallbackURL string 
	// CallbackRule represents the Rule of the callback, empty is using default rule, `all` is callback all result
	CallbackRule string 
	// RoomID represents the room id
	RoomID string 
	// UserID represents the user id
	UserID string 
	// ForumID represents the forum id
	ForumID string
}
 ```

**Perform** return a json `string` and a `int` express response, a `error` identifying success of failure

Only remote files are supported, and request information is created via structure `LongSpeech`

> **Params  Descrition**

- ***secretID***: secret-id for recognition task
- ***longspch***: LongSpeech struct wrapper long speech message for request
- ***timeout***: Set request timeout, if value equal 0, will using default timeout(30s)

---



## License

[MIT](http://www.opensource.org/licenses/mit-license.php)