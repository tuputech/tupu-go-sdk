// Package videosync provide interface of TUPU video recognition
package videosync

import (
	"fmt"

	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

// video extends recognition.DataInfo to descripton video file
type (
	VideoSync struct {
		dataInfo  *tupumodel.DataInfo
		interval  uint8
		maxFrames uint16
		tag       string
		tasks     []string
	}

	SyncOptFunc func(*VideoSync)
)

const (
	DefalutInterval  = 1
	DefaultMaxFrames = 200
)

func newVidoSync(optFuncs ...tupumodel.OptFunc) *VideoSync {
	var (
		video    = new(VideoSync)
		dataInfo = new(tupumodel.DataInfo)
	)
	dataInfo.OtherMsg = make(map[string]string)
	video.dataInfo = dataInfo
	for _, setConf := range optFuncs {
		setConf(dataInfo)
	}
	dataInfo.SetFileType("video")
	return video
}

// NewRemoteSpeech is an initializer for create Speech resource with url
func NewRemoteSpeech(url string) *VideoSync {

	// verify the params
	if tupuerror.StringIsEmpty(url) {
		return nil
	}
	video := newVidoSync(tupumodel.WithFileURL(url))
	return video
}

// NewLocalSpeech is an initializer for create Speech resource with local file path
func NewLocalSpeech(path string) *VideoSync {

	// verify the params
	if tupuerror.StringIsEmpty(path) {
		return nil
	}

	video := newVidoSync(tupumodel.WithLocalPath(path))

	return video
}

// NewBinarySpeech is an initializer for create Speech resource with binary content
func NewBinarySpeech(buf []byte, fileName string) *VideoSync {
	// verify the params
	if tupuerror.PtrIsNil(buf) || tupuerror.StringIsEmpty(fileName) {
		return nil
	}

	video := newVidoSync(tupumodel.WithBinary(buf, fileName))

	return video
}

// ClearBuffer is an helper to clear video binary content
func (video *VideoSync) ClearData() {
	video.dataInfo.ClearData()
	video.tag = ""
	video.interval = DefalutInterval
	video.maxFrames = DefaultMaxFrames
	video.tasks = nil
}

// InitConf provide uniform entry setting attributes
func (vdSync *VideoSync) InitConf(options ...tupumodel.OptFunc) {
	for _, opt := range options {
		opt(vdSync.dataInfo)
	}
}

func (vdSync *VideoSync) InitOptionParams(optFuncs ...SyncOptFunc) {
	for _, opt := range optFuncs {
		opt(vdSync)
	}
}

func WithTag(tag string) SyncOptFunc {
	return func(vs *VideoSync) {
		vs.dataInfo.OtherMsg["tag"] = tag
		vs.tag = tag
	}
}

func WithInterval(interval uint8) SyncOptFunc {
	return func(vs *VideoSync) {
		vs.dataInfo.OtherMsg["interval"] = fmt.Sprint(interval)
		vs.interval = interval
	}
}

func WithMaxFrames(maxFrames uint16) SyncOptFunc {
	return func(vs *VideoSync) {
		vs.dataInfo.OtherMsg["maxFrames"] = fmt.Sprint(maxFrames)
	}
}

func WithTask(tasks ...string) SyncOptFunc {
	return func(vs *VideoSync) {
		vs.tasks = tasks
	}
}
