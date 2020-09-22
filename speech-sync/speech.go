// Package speechsync provide interface of TUPU speech recognition
package speechsync

import (
	baseReg "github.com/tuputech/tupu-go-sdk/base-recognition"
)

// Speech extends recognition.DataInfo to descripton speech file
type Speech struct {
	dataInfo baseReg.DataInfo
}

// NewRemoteSpeech is an initializer for create Speech resource with url
func NewRemoteSpeech(url string) *Speech {
	dataInfo := new(baseReg.DataInfo)
	speech := new(Speech)

	dataInfo = baseReg.NewRemoteDataInfo(url)
	speech.dataInfo = *dataInfo
	return speech
}

// NewLocalSpeech is an initializer for create Speech resource with local file path
func NewLocalSpeech(path string) *Speech {
	dataInfo := new(baseReg.DataInfo)
	speech := new(Speech)

	dataInfo = baseReg.NewLocalDataInfo(path)
	speech.dataInfo = *dataInfo
	return speech
}

// NewBinarySpeech is an initializer for create Speech resource with binary content
func NewBinarySpeech(buf []byte, fileName string) *Speech {
	dataInfo := new(baseReg.DataInfo)
	speech := new(Speech)

	dataInfo = baseReg.NewBinaryDataInfo(buf, fileName)
	speech.dataInfo = *dataInfo
	return speech
}

// ClearBuffer is an helper to clear speech binary content
func (speech *Speech) ClearBuffer() {
	speech.dataInfo.ClearBuffer()
}
