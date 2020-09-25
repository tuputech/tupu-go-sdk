// Package shortsync provide interface of TUPU speech recognition
package shortsync

import (
	generalrcn "github.com/tuputech/tupu-go-sdk/recognition-api/general"
)

// Speech extends recognition.DataInfo to descripton speech file
type Speech struct {
	dataInfo generalrcn.DataInfo
}

// NewRemoteSpeech is an initializer for create Speech resource with url
func NewRemoteSpeech(url string) *Speech {

	// verify the params
	if generalrcn.StringIsEmpty(url) {
		return nil
	}

	var (
		dataInfo = new(generalrcn.DataInfo)
		speech   = new(Speech)
	)

	dataInfo = generalrcn.NewRemoteDataInfo(url)
	speech.dataInfo = *dataInfo
	return speech
}

// NewLocalSpeech is an initializer for create Speech resource with local file path
func NewLocalSpeech(path string) *Speech {

	// verify the params
	if generalrcn.StringIsEmpty(path) {
		return nil
	}

	var (
		dataInfo = new(generalrcn.DataInfo)
		speech   = new(Speech)
	)

	dataInfo = generalrcn.NewLocalDataInfo(path)
	speech.dataInfo = *dataInfo
	return speech
}

// NewBinarySpeech is an initializer for create Speech resource with binary content
func NewBinarySpeech(buf []byte, fileName string) *Speech {
	// verify the params
	if generalrcn.PtrIsNil(buf) || generalrcn.StringIsEmpty(fileName) {
		return nil
	}

	var (
		dataInfo = new(generalrcn.DataInfo)
		speech   = new(Speech)
	)

	dataInfo = generalrcn.NewBinaryDataInfo(buf, fileName)
	speech.dataInfo = *dataInfo
	return speech
}

// ClearBuffer is an helper to clear speech binary content
func (speech *Speech) ClearBuffer() {
	speech.dataInfo.ClearBuffer()
}
