package speechsync

import (
	basercn "github.com/tuputech/tupu-go-sdk/base-recognition"
)

const (
	SPEECH_API_URL = "http://api.open.tuputech.com/v3/recognition/speech/"
)

// SpeechHandler is a client-side helper to access TUPU speech recognition service
type SpeechHandler struct {
	hdler basercn.Handler
}

// NewSpeechHandler is an initializer for a SpeechHandler
func NewSpeechHandler(privateKeyPath string) (*SpeechHandler, error) {
	var err error
	spHdler := new(SpeechHandler)
	hdler := new(basercn.Handler)

	if hdler, err = basercn.NewHandlerWithURL(privateKeyPath, SPEECH_API_URL); err != nil {
		spHdler.hdler = *hdler
		return spHdler, err
	}
	spHdler.hdler = *hdler
	return spHdler, err
}

// NewSpeechHandlerWithURL is an initializer for a SpeechHandler with url
func NewSpeechHandlerWithURL(privateKeyPath, url string) (*SpeechHandler, error) {
	var err error
	spHdler := new(SpeechHandler)
	hdler := new(basercn.Handler)

	if hdler, err = basercn.NewHandlerWithURL(privateKeyPath, url); err != nil {
		spHdler.hdler = *hdler
		return spHdler, err
	}

	spHdler.hdler = *hdler
	return spHdler, err
}

// Perform is the major method for initiating a speech recognition request
func (spHdler *SpeechHandler) Perform(secretID string, spSlice []*Speech) (string, int, error) {

	dataInfoSlice := make([]*basercn.DataInfo, 0)
	for i := 0; i < len(spSlice); i++ {
		dataInfoSlice = append(dataInfoSlice, &spSlice[i].dataInfo)
	}
	return spHdler.hdler.Recognize("speech", secretID, dataInfoSlice, nil)
}

// PerformWithURL is a shortcut for initiating a speech recognition request with URLs
func (spHdler *SpeechHandler) PerformWithURL(secretID string, URLs []string) (string, int, error) {
	return spHdler.hdler.RecognizeWithURL("speech", secretID, URLs, nil)
}

// PerformWithPath is a shortcut for initiating a speech recognition request with paths
func (spHdler *SpeechHandler) PerformWithPath(secretID string, speechPaths []string) (string, int, error) {
	return spHdler.hdler.RecognizeWithPath("speech", secretID, speechPaths, nil)
}
