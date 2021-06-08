package speechsync

import (
	"fmt"
	"sync"

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
	syncPool sync.Pool
	hdler    *tupucontrol.Handler
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

	syncHdler.syncPool.New = func() interface{} {
		return newSpeechSync()
	}

	return syncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (syncHdler *SyncHandler) SetServerURL(url string) {
	syncHdler.hdler.SetServerURL(url)
}

// PerformWithBinary is the major method for initiating a speech recognition request, Params binaryData key is fileName(include filetype, example "1.flv"), value is binary data
func (syncHdler *SyncHandler) PerformWithBinary(secretID string, binaryData map[string][]byte, tasks ...string) (result string, statusCode int, err error) {

	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(binaryData) {
		err = fmt.Errorf("[Params Error]: caller function name: %s", tupuerror.GetCallerFuncName())
		statusCode = 400
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, len(binaryData)+1)
		index         = 0
		speechSync    *SpeechSync
	)

	// wrapper data to DataInfo
	for fileName, buf := range binaryData {
		speechSync = syncHdler.syncPool.Get().(*SpeechSync)
		defer syncHdler.recycleDataObj(speechSync)
		// set struct of dataInfo value
		speechSync.InitConf(tupumodel.WithBinary(buf, fileName))
		dataInfoSlice[index] = speechSync.dataInfo
		index++
	}
	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, tasks)
}

// PerformWithURL is a shortcut for initiating a speech recognition request with URLs
func (syncHdler *SyncHandler) PerformWithURL(secretID string, URLs []string, tasks ...string) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(URLs) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, len(URLs)+1)
		speechSync    *SpeechSync
	)

	// wrapper data to DataInfo
	for index, url := range URLs {
		speechSync = syncHdler.syncPool.Get().(*SpeechSync)
		// reset DataInfo, and put to pool
		defer syncHdler.recycleDataObj(speechSync)
		speechSync.InitConf(tupumodel.WithFileURL(url))

		dataInfoSlice[index] = speechSync.dataInfo
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, tasks)
}

// PerformWithPath is a shortcut for initiating a speech recognition request with paths
func (syncHdler *SyncHandler) PerformWithPath(secretID string, speechPaths []string, tasks ...string) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(speechPaths) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, len(speechPaths)+1)
		speechSync    *SpeechSync
	)

	// wrapper data to DataInfo
	for index, path := range speechPaths {
		// malloc and init
		speechSync = syncHdler.syncPool.Get().(*SpeechSync)
		// reset DataInfo, and put to pool
		defer syncHdler.recycleDataObj(speechSync)
		speechSync.InitConf(tupumodel.WithLocalPath(path))

		dataInfoSlice[index] = speechSync.dataInfo
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, tasks)
}

func illegalSpeechFile(fileExtend string) bool {
	switch fileExtend {
	case ".amr", ".mp3", ".wmv", ".wav", ".flv":
		return false
	default:
		return true
	}
}

func (syncHdler *SyncHandler) recycleDataObj(speechSync *SpeechSync) {
	speechSync.dataInfo.ClearData()
	syncHdler.syncPool.Put(speechSync)
}

// SetTimeout provide properties to set request ttl
func (syncHdler *SyncHandler) SetTimeout(timeout int) {
	syncHdler.hdler.SetTimeout(timeout)
}
