package shortsync

import (
	"fmt"
	"path"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

const (
	SpeechAPIURL = "http://api.open.tuputech.com/v3/recognition/speech/"
)

// ShortSpeechHandler is a client-side helper to access TUPU speech recognition service
type ShortSpeechHandler struct {
	hdler *tupucontrol.Handler
}

// NewShortSpeechHandler is an initializer for a SpeechHandler
func NewShortSpeechHandler(privateKeyPath, url string) (*ShortSpeechHandler, error) {
	// verify the params
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		err     error
		hdler   *tupucontrol.Handler
		spHdler = new(ShortSpeechHandler)
	)

	if len(url) == 0 {
		url = SpeechAPIURL
	}

	if hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, url); err != nil {
		return nil, err
	}

	spHdler.hdler = hdler
	return spHdler, nil
}

// PerformWithBinary is the major method for initiating a speech recognition request, Params binaryData key is fileName(include filetype, example "1.flv"), value is binary data
func (spHdler *ShortSpeechHandler) PerformWithBinary(secretID string, binaryData map[string][]byte, timeout int) (result string, statusCode int, err error) {

	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(binaryData) {
		err = fmt.Errorf("[Params Error]: caller function name: %s", tupuerror.GetCallerFuncName())
		statusCode = 400
		return
	}

	spHdler.hdler.SetTimeout(timeout)

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		shortSpch     *ShortSpeech
	)

	// wrapper data to DataInfo
	for fileName, buf := range binaryData {
		// verify the file extend
		extend := path.Ext(fileName)
		if illegalSpeechFile(extend) {
			err = fmt.Errorf("illegal speech file, only supports wav, wmv, mp3, flv, amr, your file is %v", extend)
			statusCode = 400
			return
		}

		shortSpch = NewBinarySpeech(buf, fileName)
		dataInfoSlice = append(dataInfoSlice, shortSpch.dataInfo)
	}
	// Do request
	return spHdler.hdler.Recognize(secretID, dataInfoSlice)
}

// PerformWithURL is a shortcut for initiating a speech recognition request with URLs
func (spHdler *ShortSpeechHandler) PerformWithURL(secretID string, URLs []string, timeout int) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(URLs) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		shortSpch     *ShortSpeech
	)

	spHdler.hdler.SetTimeout(timeout)

	// wrapper data to DataInfo
	for _, url := range URLs {
		shortSpch = NewRemoteSpeech(url)
		dataInfoSlice = append(dataInfoSlice, shortSpch.dataInfo)
	}

	// Do request
	return spHdler.hdler.Recognize(secretID, dataInfoSlice)
}

// PerformWithPath is a shortcut for initiating a speech recognition request with paths
func (spHdler *ShortSpeechHandler) PerformWithPath(secretID string, speechPaths []string, timeout int) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(speechPaths) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		shortSpch     *ShortSpeech
	)

	spHdler.hdler.SetTimeout(timeout)

	// wrapper data to DataInfo
	for _, path := range speechPaths {
		shortSpch = NewLocalSpeech(path)
		dataInfoSlice = append(dataInfoSlice, shortSpch.dataInfo)
	}

	// Do request
	return spHdler.hdler.Recognize(secretID, dataInfoSlice)
}

func illegalSpeechFile(fileExtend string) bool {
	switch fileExtend {
	case ".amr", ".mp3", ".wmv", ".wav", ".flv":
		return false
	default:
		return true
	}
}
