// Package textsync provide interface of TUPU text sync recognition
package textsync

import (
	"encoding/json"
	"fmt"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

const (
	// TextSyncAPIURL is default
	TextSyncAPIURL = "http://api.open.tuputech.com/v3/recognition/text/"
)

// SyncHandler is a client-side helper to access TUPU sync text recognition service
type SyncHandler struct {
	hdler *tupucontrol.Handler
}

// NewTextHandler is an initializer for a SyncHandler.
func NewTextHandler(privateKeyPath string) (*SyncHandler, error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(privateKeyPath) {
		return nil, fmt.Errorf("[Params ERROR]: function name is %s", tupuerror.GetCallerFuncName())
	}

	var (
		err        error
		asyncHdler = new(SyncHandler)
	)

	// create TUPU general Handler
	if asyncHdler.hdler, err = tupucontrol.NewHandlerWithURL(privateKeyPath, TextSyncAPIURL); err != nil {
		return nil, err
	}

	return asyncHdler, nil
}

// SetServerURL provide set request server URL attribute
func (asyncHdler *SyncHandler) SetServerURL(url string) {
	asyncHdler.hdler.SetServerURL(url)
}

// Perform is the major method for initiating a text recognition request
func (asyncHdler *SyncHandler) Perform(secretID string, textSync []TextAsyncItem) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerror.StringIsEmpty(secretID) || tupuerror.PtrIsNil(textSync) {
		statusCode = 400
		err = fmt.Errorf("[Params ERROR]: now func: %s\tcaller func: %s ", tupuerror.GetCurrentFuncName(), tupuerror.GetCallerFuncName())
		return
	}

	var (
		recording     []byte
		requestParams string
	)

	// step2. serialize to JSON string
	recording, _ = json.Marshal(textSync)
	requestParams = string(recording)
	requestParams = `"text":` + requestParams

	// step3. transfer general api
	return asyncHdler.hdler.RecognizeWithJSON(requestParams, secretID)
}

// SetTimeout provide properties to set request ttl
func (asyncHdler *SyncHandler) SetTimeout(timeout int) {
	asyncHdler.hdler.SetTimeout(timeout)
}
