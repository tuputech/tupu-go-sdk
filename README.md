# TUPU GO SDK

Golang SDK for TUPU visual recognition service (v1.6.1)
######  
<https://www.tuputech.com>

## Changelogs
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
go get github.com/tuputech/tupu-go-sdk/recognition

## Example

1. [Image recognition interface example](./example/image.go)  
2. [shortSpeech recognition interface example](./example/short-speech.go)  
3. [longSpeech recognition interface example](./example/long-speech.go)  

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
