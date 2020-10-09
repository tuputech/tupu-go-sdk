# Package speech-async

## Installation
go get github.com/tuputech/tupu-go-sdk/recognition

## Example
For the returned json format string, the corresponding structure analysis is no longer provided, you can use a package similar to simplejson to parse the json string and get the fields you want.

[longSpeech recognition interface example](./example/long-speech.go)

----------------------

### func Perform
func (spHdler *LongSpeechHandler) Perform(secretID string, longspch *LongSpeech, timeout int) (result string, statusCode int, err error)

- **secretID**: secret-id for recognition tasks
- **longspch**: LongSpeech struct wrapper long speech message for request
- **timeout**: Set request timeout, if value equal 0, will using default timeout(30s)

----------------------