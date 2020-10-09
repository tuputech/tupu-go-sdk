# Package speech-sync

> Suitable for TUPU short speech recognition interface, providing access


## Example
For the returned json format string, the corresponding structure analysis is no longer provided, you can use a package similar to simplejson to parse the json string and get the fields you want.

> The processing entry is `ShortSpeechHandler` struct, using its three methods(`PerformWithBinary()`, `PerformWithURL()`, `PerformWithPath()`) to get the recognition results
> [shortSpeech recognition interface example](github.com/tuputech/tupu-go-sdk/example/short-speech.go)

   1. `func (spHdler *ShortSpeechHandler) PerformWithBinary(secretID string, binaryData map[string][]byte, timeout int) (result string, statusCode int, err error)`
      - **params**
      - 1. *secretID*: secret-id for recognition tasks
      - 2. *binaryData*: map type, key means file name, value means binary data
      - 3. *timeout*: Set request timeout, if value equal 0, will using default timeout(30s)

   2. `func (spHdler *ShortSpeechHandler) PerformWithPath(secretID string, speechPaths []string, timeout int) (result string, statusCode int, err error)`
      - **params**
      - 1. *secretID*: secret-id for recognition tasks
      - 2. *speechPaths*: local short speech paths
      - 3. *timeout*: Set request timeout, if value equal 0, will using default timeout(30s)

   3. `func (spHdler *ShortSpeechHandler) PerformWithURL(secretID string, URLs []string, timeout int) (result string, statusCode int, err error)`
      - **params**
      - 1. *secretID*: secret-id for recognition tasks
      - 2. *URLs*: remote short speech address
      - 3. *timeout*: Set request timeout, if value equal 0, will using default timeout(30s)