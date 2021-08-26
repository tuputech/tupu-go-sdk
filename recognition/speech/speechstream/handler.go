package speechstream

import (
	"encoding/json"
	"fmt"
	"sync"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

const (
	// SpeechStreamURL is default recording service address
	SpeechStreamURL       = "http://api.open.tuputech.com/v3/recognition/speech/stream/"
	SpeechStreamCloseURL  = "http://api.open.tuputech.com/v3/recognition/speech/stream/close/"
	SpeechStreamSearchURL = "http://api.open.tuputech.com/v3/recognition/speech/stream/search/"
)

// SpeechStreamHandler is a client-side helper to access TUPU speech recognition service
type SpeechStreamHandler struct {
	syncPool sync.Pool
	hdler    *tupucontrol.Handler
}

// NewASyncHandler is an initializer for a SpeechHandler
func NewSpeechStreamHandler(privateKeyPath string) (*SpeechStreamHandler, error) {
	// verify the params
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		err         error
		spstrmHdler = new(SpeechStreamHandler)
	)

	if spstrmHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, SpeechStreamURL); err != nil {
		return nil, err
	}

	spstrmHdler.syncPool.New = func() interface{} {
		return newSpeechStream()
	}

	return spstrmHdler, nil
}

// SetServerURL provide set request server URL attribute
func (spstrmHdler *SpeechStreamHandler) SetServerURL(url string) {
	spstrmHdler.hdler.SetServerURL(url)
}

func (spstrmHdler *SpeechStreamHandler) recycleDataObj(speechStream *SpeechStream) {
	speechStream.ClearData()
	spstrmHdler.syncPool.Put(speechStream)
}

// SetTimeout provide properties to set request ttl
func (spstrmHdler *SpeechStreamHandler) SetTimeout(timeout int) {
	spstrmHdler.hdler.SetTimeout(timeout)
}

// Perform is the major method for initiating a recognition request
func (spstrmHdler *SpeechStreamHandler) StartStreamRecognition(secretID, streamUrl, callbackUrl string, optFuncs ...StreamOptFunc) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(secretID, streamUrl, callbackUrl) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	var (
		speechStream  *SpeechStream
		requestParams []byte
		paramsStr     string
	)

	spstrmHdler.hdler.SetServerURL(SpeechStreamURL)
	speechStream = spstrmHdler.syncPool.Get().(*SpeechStream)
	defer spstrmHdler.recycleDataObj(speechStream)

	// set optional params
	speechStream.URL = streamUrl
	speechStream.Callback = callbackUrl

	for _, setConf := range optFuncs {
		setConf(speechStream)
	}

	// step2. serialize to JSON string
	requestParams, _ = json.Marshal(speechStream)

	paramsStr = `"speechStream":[` + string(requestParams) + `]`
	if speechStream.tasks != nil {
		taskStrSlice, _ := json.Marshal(speechStream.tasks)
		paramsStr += `,"tasks":` + string(taskStrSlice)
	}
	// step3. transfer general api
	return spstrmHdler.hdler.RecognizeWithJSON(paramsStr, secretID)
}

// CloseRecognitionTask can close your speech recognition task by requestId
func (spstrmHdler *SpeechStreamHandler) CloseRecognitionTask(secretID, requestId string) (result string, statusCode int, err error) {
	if tupuerror.StringIsEmpty(secretID, requestId) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}
	spstrmHdler.SetServerURL(SpeechStreamCloseURL)
	requestParams := `"speechStream":[{"requestId": "` + requestId + `"}]`
	return spstrmHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}

// QueryStatus can query your video recognition result by requestId
func (spstrmHdler *SpeechStreamHandler) QueryStatus(secretID, requestId string) (result string, statusCode int, err error) {
	if tupuerror.StringIsEmpty(secretID, requestId) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}
	spstrmHdler.SetServerURL(SpeechStreamSearchURL)

	requestParams := `"requestId": "` + requestId + `"`
	return spstrmHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}
