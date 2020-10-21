package longasync

import (
	"encoding/json"
	"fmt"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

const (
	LongSpeechAPIURL = "http://api.open.tuputech.com/v3/recognition/speech/recording/async/"
)

// LongSpeechHandler is a client-side helper to access TUPU long speech recognition service
type LongSpeechHandler struct {
	hdler *tupucontrol.Handler
}

// NewSpeechHandler is an initializer for a LongSpeechHandler. If url-param is empty, the default url is used
func NewSpeechHandler(privateKeyPath, url string) (*LongSpeechHandler, error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("[Params ERROR]: function name is %s", tupuerror.GetCallerFuncName())
	}

	var (
		err     error
		spHdler = new(LongSpeechHandler)
		hdler   = new(tupucontrol.Handler)
	)

	// using caller url or default url
	if len(url) == 0 {
		url = LongSpeechAPIURL
	}
	// create TUPU general Handler
	if hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, LongSpeechAPIURL); err != nil {
		return nil, err
	}

	spHdler.hdler = hdler
	return spHdler, nil
}

// Perform is the major method for initiating a speech recognition request
func (spHdler *LongSpeechHandler) Perform(secretID string, longspch *LongSpeech, timeout int) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(longspch) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	var (
		recording     []byte
		requestParams string
	)

	if timeout != 0 {
		spHdler.hdler.SetTimeout(timeout)
	}

	// step2. serialize to JSON string
	recording, _ = json.Marshal(longspch)
	requestParams = string(recording)
	requestParams = `"recording":` + string(recording)

	// step3. transfer general api
	return spHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}
