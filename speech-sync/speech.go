// Package speechsync provide interface of TUPU speech recognition
package speechsync

import (
	basercn "github.com/tuputech/tupu-go-sdk/base-recognition"
)

// Speech extends recognition.DataInfo to descripton speech file
type Speech struct {
	dataInfo basercn.DataInfo
}

// NewRemoteSpeech is an initializer for create Speech resource with url
func NewRemoteSpeech(url string) *Speech {
	dataInfo := new(basercn.DataInfo)
	speech := new(Speech)

	dataInfo = basercn.NewRemoteDataInfo(url)
	speech.dataInfo = *dataInfo
	return speech
}

// NewLocalSpeech is an initializer for create Speech resource with local file path
func NewLocalSpeech(path string) *Speech {
	dataInfo := new(basercn.DataInfo)
	speech := new(Speech)

	dataInfo = basercn.NewLocalDataInfo(path)
	speech.dataInfo = *dataInfo
	return speech
}

// NewBinarySpeech is an initializer for create Speech resource with binary content
func NewBinarySpeech(buf []byte, fileName string) *Speech {
	dataInfo := new(basercn.DataInfo)
	speech := new(Speech)

	dataInfo = basercn.NewBinaryDataInfo(buf, fileName)
	speech.dataInfo = *dataInfo
	return speech
}

// ClearBuffer is an helper to clear speech binary content
func (speech *Speech) ClearBuffer() {
	speech.dataInfo.ClearBuffer()
}
