package videoasync

import (
	"encoding/json"
	"fmt"
	"sync"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

const (
	// VideoSyncURL is default recording service address
	VideoAsyncURL          = "http://api.open.tuputech.com/v3/recognition/video/asyncscan/"
	VideoAsyncResultURL    = "http://api.open.tuputech.com/v3/recognition/video/result/"
	VideoAsyncCloseTaskURL = "http://api.open.tuputech.com/v3/recognition/video/close/"
	VideoAsyncQueryRateURL = "http://api.open.tuputech.com/v3/recognition/video/rate/"
)

// AsyncHandler is a client-side helper to access TUPU speech recognition service
type AsyncHandler struct {
	syncPool sync.Pool
	hdler    *tupucontrol.Handler
}

// NewASyncHandler is an initializer for a SpeechHandler
func NewVideoAsyncHandler(privateKeyPath string) (*AsyncHandler, error) {
	// verify the params
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		err        error
		asyncHdler = new(AsyncHandler)
	)

	if asyncHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, VideoAsyncURL); err != nil {
		return nil, err
	}

	asyncHdler.syncPool.New = func() interface{} {
		return newVidoASync()
	}

	return asyncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (asyncHdler *AsyncHandler) SetServerURL(url string) {
	asyncHdler.hdler.SetServerURL(url)
}

func (syncHdler *AsyncHandler) recycleDataObj(videoAsync *VideoAsync) {
	videoAsync.ClearData()
	syncHdler.syncPool.Put(videoAsync)
}

// SetTimeout provide properties to set request ttl
func (asyncHdler *AsyncHandler) SetTimeout(timeout int) {
	asyncHdler.hdler.SetTimeout(timeout)
}

// Perform is the major method for initiating a recognition request
func (asyncHdler *AsyncHandler) Perform(secretID, videoUrl, callbackUrl string, optFuncs ...AsyncOptFunc) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(secretID, videoUrl, callbackUrl) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	var (
		videoAsync    *VideoAsync
		requestParams []byte
		paramsStr     string
	)

	asyncHdler.hdler.SetServerURL(VideoAsyncURL)
	videoAsync = asyncHdler.syncPool.Get().(*VideoAsync)
	defer asyncHdler.recycleDataObj(videoAsync)

	// set optional params
	videoAsync.Video = videoUrl
	videoAsync.CallbackUrl = callbackUrl

	for _, setConf := range optFuncs {
		setConf(videoAsync)
	}

	// step2. serialize to JSON string
	requestParams, _ = json.Marshal(videoAsync)

	paramsStr = string(requestParams[1 : len(requestParams)-1])
	// step3. transfer general api
	return asyncHdler.hdler.RecognizeWithJSON(paramsStr, secretID)
}

// CloseRecognitionTask can close your video recognition task
func (asyncHdler *AsyncHandler) CloseRecognitionTask(secretID, videoId string) (result string, statusCode int, err error) {
	asyncHdler.SetServerURL(VideoAsyncCloseTaskURL)
	return asyncHdler.closeOrQueryVideoInfo(secretID, videoId)
}

// QueryRecognitionResult can query your video recognition result
func (asyncHdler *AsyncHandler) QueryRecognitionResult(secretID, videoId string) (result string, statusCode int, err error) {
	asyncHdler.SetServerURL(VideoAsyncResultURL)
	return asyncHdler.closeOrQueryVideoInfo(secretID, videoId)
}

// QueryRecognitionResult can query video recognition rate for your secretId
func (asyncHdler *AsyncHandler) QueryRate(secretID string) (result string, statusCode int, err error) {
	if tupuerror.StringIsEmpty(secretID) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}
	asyncHdler.SetServerURL(VideoAsyncQueryRateURL)

	return asyncHdler.hdler.RecognizeWithJSON("{}", secretID)
}

func (asyncHdler *AsyncHandler) closeOrQueryVideoInfo(secretID, videoId string) (result string, statusCode int, err error) {
	if tupuerror.StringIsEmpty(secretID, videoId) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	requestParams := `"videoId": "` + videoId + `"`
	return asyncHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}
