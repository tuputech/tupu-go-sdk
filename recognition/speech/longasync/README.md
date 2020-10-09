# Package speech-async

## Installation
go get github.com/tuputech/tupu-go-sdk/recognition

## Example

[longSpeech recognition interface example](./example/long-speech.go)

----------------------

### func Perform
func (spHdler *LongSpeechHandler) Perform(secretID string, longspch *LongSpeech, timeout int) (result string, statusCode int, err error)

- **secretID**: secret-id for recognition tasks
- **longspch**: LongSpeech struct wrapper long speech message for request
- **timeout**: setting http request timeout

----------------------