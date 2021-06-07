// Package speechsync provide interface of TUPU speech recognition
package speechsync

import (
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

// SpeechSync extends recognition.DataInfo to descripton speech file
type (
	SpeechSync struct {
		dataInfo *tupumodel.DataInfo
	}

	SyncOptFunc func(*SpeechSync)
)

func newSpeechSync(optFuncs ...tupumodel.OptFunc) *SpeechSync {
	var (
		speech   = new(SpeechSync)
		dataInfo = new(tupumodel.DataInfo)
	)
	dataInfo.SetFileType("speech")
	speech.dataInfo = dataInfo
	for _, setConf := range optFuncs {
		setConf(dataInfo)
	}
	return speech
}

// NewRemoteSpeech is an initializer for create Speech resource with url
func NewRemoteSpeech(url string) *SpeechSync {

	// verify the params
	if tupuerror.StringIsEmpty(url) {
		return nil
	}
	speech := newSpeechSync()
	speech.dataInfo.SetRemoteInfo(url)
	return speech
}

// NewLocalSpeech is an initializer for create Speech resource with local file path
func NewLocalSpeech(path string) *SpeechSync {

	// verify the params
	if tupuerror.StringIsEmpty(path) {
		return nil
	}

	speech := newSpeechSync()
	speech.dataInfo.SetPath(path)

	return speech
}

// NewBinarySpeech is an initializer for create Speech resource with binary content
func NewBinarySpeech(buf []byte, fileName string) *SpeechSync {
	// verify the params
	if tupuerror.PtrIsNil(buf) || tupuerror.StringIsEmpty(fileName) {
		return nil
	}

	speech := newSpeechSync()
	speech.dataInfo.SetBuf(buf)
	speech.dataInfo.SetFileName(fileName)

	return speech
}

// ClearBuffer is an helper to clear speech binary content
func (speech *SpeechSync) ClearBuffer() {
	speech.dataInfo.ClearBuffer()
}

// InitConf provide uniform entry setting attributes
func (spSync *SpeechSync) InitConf(options ...tupumodel.OptFunc) {
	for _, opt := range options {
		opt(spSync.dataInfo)
	}
}
