// Package speechasync provide interface of TUPU speech async recognition
package speechasync

import (
	"encoding/json"
	"fmt"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

const (
	// SpeechAsyncAPIURL is default
	SpeechAsyncAPIURL = "http://api.open.tuputech.com/v3/recognition/speech/recording/async/"
)

// AsyncHandler is a client-side helper to access TUPU async speech recognition service
type AsyncHandler struct {
	hdler *tupucontrol.Handler
}

// NewSpeechHandler is an initializer for a AsyncHandler. If url-param is empty, the default url is used
func NewSpeechHandler(privateKeyPath string) (*AsyncHandler, error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("[Params ERROR]: function name is %s", tupuerror.GetCallerFuncName())
	}

	var (
		err        error
		asyncHdler = new(AsyncHandler)
	)

	// create TUPU general Handler
	if asyncHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, SpeechAsyncAPIURL); err != nil {
		return nil, err
	}

	return asyncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (asyncHdler *AsyncHandler) SetServerURL(url string) {
	asyncHdler.hdler.SetServerURL(url)
}

// Perform is the major method for initiating a speech recognition request
func (asyncHdler *AsyncHandler) Perform(secretID string, speechAsync *SpeechAsync) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(speechAsync) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	var (
		recording     []byte
		requestParams string
	)

	// step2. serialize to JSON string
	recording, _ = json.Marshal(speechAsync)
	requestParams = string(recording)
	requestParams = `"recording":` + string(recording)

	// step3. transfer general api
	return asyncHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}

// SetTimeout provide properties to set request ttl
func (asyncHdler *AsyncHandler) SetTimeout(timeout int) {
	asyncHdler.hdler.SetTimeout(timeout)
}
