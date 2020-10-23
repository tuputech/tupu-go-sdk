package speechsync

import (
	"fmt"
	"path"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

const (
	// SpeechSyncAPIURL is default recording service address
	SpeechSyncAPIURL = "http://api.open.tuputech.com/v3/recognition/speech/"
)

// SyncHandler is a client-side helper to access TUPU speech recognition service
type SyncHandler struct {
	hdler *tupucontrol.Handler
}

// NewSyncHandler is an initializer for a SpeechHandler
func NewSyncHandler(privateKeyPath string) (*SyncHandler, error) {
	// verify the params
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		err       error
		syncHdler = new(SyncHandler)
	)

	if syncHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, SpeechSyncAPIURL); err != nil {
		return nil, err
	}

	return syncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (syncHdler *SyncHandler) SetServerURL(url string) {
	syncHdler.hdler.SetServerURL(url)
}

// PerformWithBinary is the major method for initiating a speech recognition request, Params binaryData key is fileName(include filetype, example "1.flv"), value is binary data
func (syncHdler *SyncHandler) PerformWithBinary(secretID string, binaryData map[string][]byte) (result string, statusCode int, err error) {

	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(binaryData) {
		err = fmt.Errorf("[Params Error]: caller function name: %s", tupuerror.GetCallerFuncName())
		statusCode = 400
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		speechSync    *SpeechSync
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

		speechSync = NewBinarySpeech(buf, fileName)
		dataInfoSlice = append(dataInfoSlice, speechSync.dataInfo)
	}
	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice)
}

// PerformWithURL is a shortcut for initiating a speech recognition request with URLs
func (syncHdler *SyncHandler) PerformWithURL(secretID string, URLs []string) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(URLs) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		speechSync    *SpeechSync
	)

	// wrapper data to DataInfo
	for _, url := range URLs {
		speechSync = NewRemoteSpeech(url)
		dataInfoSlice = append(dataInfoSlice, speechSync.dataInfo)
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice)
}

// PerformWithPath is a shortcut for initiating a speech recognition request with paths
func (syncHdler *SyncHandler) PerformWithPath(secretID string, speechPaths []string) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(speechPaths) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
		speechSync    *SpeechSync
	)

	// wrapper data to DataInfo
	for _, path := range speechPaths {
		speechSync = NewLocalSpeech(path)
		dataInfoSlice = append(dataInfoSlice, speechSync.dataInfo)
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice)
}

func illegalSpeechFile(fileExtend string) bool {
	switch fileExtend {
	case ".amr", ".mp3", ".wmv", ".wav", ".flv":
		return false
	default:
		return true
	}
}

// SetTimeout provide properties to set request ttl
func (syncHdler *SyncHandler) SetTimeout(timeout int) {
	syncHdler.hdler.SetTimeout(timeout)
}
