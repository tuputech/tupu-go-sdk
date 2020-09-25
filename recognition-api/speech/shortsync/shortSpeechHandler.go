package shortsync

import (
	"fmt"

	generalrcn "github.com/tuputech/tupu-go-sdk/recognition-api/general"
)

const (
	SPEECH_API_URL = "http://api.open.tuputech.com/v3/recognition/speech/"
)

// SpeechHandler is a client-side helper to access TUPU speech recognition service
type SpeechHandler struct {
	hdler generalrcn.Handler
}

// NewSpeechHandler is an initializer for a SpeechHandler
func NewSpeechHandler(privateKeyPath string) (*SpeechHandler, error) {
	// verify the params
	if generalrcn.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("[Params Error]: caller function name: %s", generalrcn.GetCallerFuncName())
	}
	var (
		err     error
		spHdler = new(SpeechHandler)
		hdler   = new(generalrcn.Handler)
	)

	if hdler, err = generalrcn.NewHandlerWithURL(privateKeyPath, SPEECH_API_URL); err != nil {
		spHdler.hdler = *hdler
		return spHdler, err
	}
	spHdler.hdler = *hdler
	return spHdler, err
}

// NewSpeechHandlerWithURL is an initializer for a SpeechHandler with url
func NewSpeechHandlerWithURL(privateKeyPath, url string) (*SpeechHandler, error) {
	// verify the params
	if generalrcn.StringIsEmpty(privateKeyPath, url) {
		return nil, fmt.Errorf("[Params Error]: caller function name: %s", generalrcn.GetCallerFuncName())
	}
	var (
		err     error
		spHdler = new(SpeechHandler)
		hdler   = new(generalrcn.Handler)
	)

	if hdler, err = generalrcn.NewHandlerWithURL(privateKeyPath, url); err != nil {
		spHdler.hdler = *hdler
		return spHdler, err
	}

	spHdler.hdler = *hdler
	return spHdler, err
}

// Perform is the major method for initiating a speech recognition request
func (spHdler *SpeechHandler) Perform(secretID string, spSlice []*Speech) (string, int, error) {

	// verify the params
	if generalrcn.StringIsEmpty(secretID) || generalrcn.PtrIsNil(spSlice) {
		return "", 400, fmt.Errorf("[Params Error]: caller function name: %s", generalrcn.GetCallerFuncName())
	}

	dataInfoSlice := make([]*generalrcn.DataInfo, 0)

	for i := 0; i < len(spSlice); i++ {
		dataInfoSlice = append(dataInfoSlice, &spSlice[i].dataInfo)
	}
	return spHdler.hdler.Recognize("speech", secretID, dataInfoSlice, nil)
}

// PerformWithURL is a shortcut for initiating a speech recognition request with URLs
func (spHdler *SpeechHandler) PerformWithURL(secretID string, URLs []string) (string, int, error) {
	// verify the params
	if generalrcn.StringIsEmpty(secretID) || generalrcn.PtrIsNil(URLs) {
		return "", 400, fmt.Errorf("[Params Error]: caller function name: %s", generalrcn.GetCallerFuncName())
	}
	return spHdler.hdler.RecognizeWithURL("speech", secretID, URLs, nil)
}

// PerformWithPath is a shortcut for initiating a speech recognition request with paths
func (spHdler *SpeechHandler) PerformWithPath(secretID string, speechPaths []string) (string, int, error) {
	// verify the params
	if generalrcn.StringIsEmpty(secretID) || generalrcn.PtrIsNil(speechPaths) {
		return "", 400, fmt.Errorf("[Params Error]: caller function name: %s", generalrcn.GetCallerFuncName())
	}
	return spHdler.hdler.RecognizeWithPath("speech", secretID, speechPaths, nil)
}
