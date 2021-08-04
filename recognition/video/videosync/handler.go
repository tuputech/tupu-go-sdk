package videosync

import (
	"fmt"
	"sync"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

const (
	// VideoSyncURL is default recording service address
	VideoSyncURL = "http://api.open.tuputech.com/v3/recognition/video/syncscan/"
)

// SyncHandler is a client-side helper to access TUPU video recognition service
type SyncHandler struct {
	syncPool sync.Pool
	hdler    *tupucontrol.Handler
}

// NewSyncHandler is an initializer for a VideoHandler
func NewSyncHandler(privateKeyPath string) (*SyncHandler, error) {
	// verify the params
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		err       error
		syncHdler = new(SyncHandler)
	)

	if syncHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, VideoSyncURL); err != nil {
		return nil, err
	}

	syncHdler.syncPool.New = func() interface{} {
		return newVidoSync()
	}

	return syncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (syncHdler *SyncHandler) SetServerURL(url string) {
	syncHdler.hdler.SetServerURL(url)
}

// PerformWithBinary is the major method for initiating a video recognition request, Params binaryData key is fileName(include filetype, example "1.flv"), value is binary data
func (syncHdler *SyncHandler) PerformWithBinary(secretID string, binaryData map[string][]byte, optFuncs ...SyncOptFunc) (result string, statusCode int, err error) {

	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(binaryData) {
		err = fmt.Errorf("[Params Error]: caller function name: %s", tupuerror.GetCallerFuncName())
		statusCode = 400
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0, len(binaryData))
		videoSync     *VideoSync
	)

	// wrapper data to DataInfo
	for fileName, buf := range binaryData {
		videoSync = syncHdler.syncPool.Get().(*VideoSync)
		defer syncHdler.recycleDataObj(videoSync)
		// set struct of dataInfo value
		videoSync.InitConf(tupumodel.WithBinary(buf, fileName))
		videoSync.InitOptionParams(optFuncs...)
		dataInfoSlice = append(dataInfoSlice, videoSync.dataInfo)
	}
	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, videoSync.tasks)
}

// PerformWithURL is a shortcut for initiating a video recognition request with URLs
func (syncHdler *SyncHandler) PerformWithURL(secretID string, URLs []string, optfuncs ...SyncOptFunc) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(URLs) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0, len(URLs))
		videoSync     *VideoSync
	)

	// wrapper data to DataInfo
	for _, url := range URLs {
		videoSync = syncHdler.syncPool.Get().(*VideoSync)
		// reset DataInfo, and put to pool
		defer syncHdler.recycleDataObj(videoSync)
		videoSync.InitConf(tupumodel.WithFileURL(url))
		videoSync.InitOptionParams(optfuncs...)

		dataInfoSlice = append(dataInfoSlice, videoSync.dataInfo)
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, videoSync.tasks)
}

// PerformWithPath is a shortcut for initiating a video recognition request with paths
func (syncHdler *SyncHandler) PerformWithPath(secretID string, speechPaths []string, optFuncs ...SyncOptFunc) (result string, statusCode int, err error) {
	// verify the params
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(speechPaths) {
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
		return
	}

	var (
		dataInfoSlice = make([]*tupumodel.DataInfo, 0, len(speechPaths))
		videoSync     *VideoSync
	)

	// wrapper data to DataInfo
	for _, path := range speechPaths {
		// malloc and init
		videoSync = syncHdler.syncPool.Get().(*VideoSync)
		// reset DataInfo, and put to pool
		defer syncHdler.recycleDataObj(videoSync)
		videoSync.InitConf(tupumodel.WithLocalPath(path))
		videoSync.InitOptionParams(optFuncs...)

		// TODO
		dataInfoSlice = append(dataInfoSlice, videoSync.dataInfo)
	}

	// Do request
	return syncHdler.hdler.Recognize(secretID, dataInfoSlice, videoSync.tasks)
}

func illegalSpeechFile(fileExtend string) bool {
	switch fileExtend {
	case ".mkv", ".mp4", ".wmv", ".rmvb", ".flv", ".3gp", ".ts", ".mov", ".gif", ".m3u8", ".mpg/mpeg", ".mxf":
		return false
	default:
		return true
	}
}

func (syncHdler *SyncHandler) recycleDataObj(videoSync *VideoSync) {
	videoSync.ClearData()
	syncHdler.syncPool.Put(videoSync)
}

// SetTimeout provide properties to set request ttl
func (syncHdler *SyncHandler) SetTimeout(timeout int) {
	syncHdler.hdler.SetTimeout(timeout)
}
