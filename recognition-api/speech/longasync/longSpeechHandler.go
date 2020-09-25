package longasync

import (
	"encoding/json"
	"fmt"

	generalrcn "github.com/tuputech/tupu-go-sdk/recognition-api/general"
)

const (
	LONG_SPEECH_API_URL = "http://api.open.tuputech.com/v3/recognition/speech/recording/async/"
)

// LongSpeechHandler is a client-side helper to access TUPU long speech recognition service
type LongSpeechHandler struct {
	generalrcn.Handler
}

// NewSpeechHandler is an initializer for a LongSpeechHandler
func NewSpeechHandler(privateKeyPath string) (*LongSpeechHandler, error) {

	// step1. Invalid parameter check
	if generalrcn.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("[Params ERROR]: function name is %s", generalrcn.GetCallerFuncName())
	}

	var (
		err     error
		spHdler = new(LongSpeechHandler)
		hdler   = new(generalrcn.Handler)
	)

	if hdler, err = generalrcn.NewHandlerWithURL(privateKeyPath, LONG_SPEECH_API_URL); err != nil {
		spHdler.Handler = *hdler
		return nil, err
	}
	spHdler.Handler = *hdler
	return spHdler, err
}

// NewSpeechHandlerWithURL is an initializer for a LongSpeechHandler with url
func NewSpeechHandlerWithURL(privateKeyPath, url string) (*LongSpeechHandler, error) {

	// step1. Invalid parameter check
	if generalrcn.StringIsEmpty(privateKeyPath, url) {
		return nil, fmt.Errorf("[Params ERROR]: function name is %s", generalrcn.GetCallerFuncName())
	}
	var (
		err     error
		spHdler = new(LongSpeechHandler)
		hdler   = new(generalrcn.Handler)
	)

	if hdler, err = generalrcn.NewHandlerWithURL(privateKeyPath, url); err != nil {
		spHdler.Handler = *hdler
		return nil, err
	}
	spHdler.Handler = *hdler
	return spHdler, err
}

// Perform is the major method for initiating a speech recognition request
func (spHdler *LongSpeechHandler) Perform(secretID string, dataType interface{}) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if generalrcn.StringIsEmpty(secretID) || generalrcn.PtrIsNil(dataType) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", generalrcn.GetCurrentFuncName(), generalrcn.GetCallerFuncName())
		return
	}

	var (
		recording     []byte
		requestParams string
	)

	// step2. serialize to JSON string
	switch longSpch := dataType.(type) {
	case *LongSpeech, LongSpeech:
		recording, _ = json.Marshal(longSpch)
		requestParams = string(recording)
		requestParams = `"recording":` + string(recording)
	default:
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", generalrcn.GetCurrentFuncName(), generalrcn.GetCallerFuncName())
		return
	}

	// step3. transfer general api
	return spHdler.RecognizeWithJSON(requestParams, secretID)
}
