// Package controller provide General functions of TUPU content recognition interface
package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tupuerrorlib "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
	tuputools "github.com/tuputech/tupu-go-sdk/lib/tools"
)

var (
	// ErrBadGateway is the error of request
	ErrBadGateway = errors.New("502 Bad Gateway")
	// ErrServiceUnavailable is the error of request
	ErrServiceUnavailable = errors.New("503 Service Unavailable")
)

const (
	// RootAPIURL is default entry of the TUPU recognition API
	RootAPIURL = "http://api.open.tuputech.com/v3/recognition/"
	// DefaultUserAgent is default value of the request Header: User-Agent
	DefaultUserAgent = "tupu-client/1.0"
	// DefaultContentType is default value of the request Header: Content-Type
	DefaultContentType = "multipart/form-data"
)

// Handler is a client-side helper to access TUPU recognition service
type Handler struct {
	Client   *http.Client
	apiURL   string
	signer   tuputools.Signer
	verifier tuputools.Verifier
	//for sub-user statistics and billing
	UID string
	// UserAgent is the request Header: User-Agent
	UserAgent string
	// ContentType is the request Header: Content-Type
	ContentType string
	// Timeout is the request Header: Timeout
	Timeout string
	// Client is the *http.Client object
}

// NewHandlerWithURL is also an initializer for a Handler
func NewHandlerWithURL(privateKeyPath, url string) (hdler *Handler, e error) {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(privateKeyPath, url) {
		return nil, fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}
	hdler = new(Handler)

	// init other default proprety
	hdler.initHandler()
	hdler.apiURL = RootAPIURL
	hdler.apiURL = url

	if hdler.verifier, e = tuputools.LoadTupuPublicKey(); e != nil {
		return nil, e
	}
	if hdler.signer, e = tuputools.LoadPrivateKey(privateKeyPath); e != nil {
		return nil, e
	}
	return hdler, nil
}

// SetTimeout is the Handler method to setting the UserAgent attribute
func (hdler *Handler) SetTimeout(timeout int) {
	if timeout != 0 {
		hdler.Timeout = fmt.Sprintf("%d", timeout)
	}
}

// SetServerURL provide setting server URL attribute
func (hdler *Handler) SetServerURL(url string) {
	hdler.apiURL = url
}

// SetContentType is the Handler method to setting the UserAgent attribute
func (hdler *Handler) SetContentType(contentType string) {
	if tupuerrorlib.StringIsEmpty(contentType) {
		return
	}
	hdler.ContentType = contentType
}

// SetUserAgent is the Handler method to setting the UserAgent attribute
func (hdler *Handler) SetUserAgent(userAgent string) {
	if tupuerrorlib.StringIsEmpty(userAgent) {
		return
	}
	hdler.UserAgent = userAgent
}

// SetUID is the Handler method to setting the UID attribute
func (hdler *Handler) SetUID(uid string) {
	if tupuerrorlib.StringIsEmpty(uid) {
		return
	}
	hdler.UID = uid
}

func (hdler *Handler) initHandler() {
	hdler.UserAgent = DefaultUserAgent
	hdler.ContentType = DefaultContentType
	hdler.Timeout = "30"
	hdler.Client = &http.Client{}
}

// RecognizeWithJSON is one of major method to access recognition api
func (hdler *Handler) RecognizeWithJSON(jsonStr, secretID string) (result string, statusCode int, err error) {

	// step1. Invalid parameter check
	if tupuerrorlib.StringIsEmpty(jsonStr, secretID) {
		result = ""
		statusCode = 400
		err = fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
		return
	}

	var (
		params    map[string]string
		paramsStr string
		url       = hdler.apiURL + secretID
		req       *http.Request
		resp      *http.Response
	)

	// step2. get timestamp, nonce, signature
	if params, err = hdler.GetGeneralParams(secretID); err != nil {
		statusCode = 400
		return
	}

	// step3. serialize to JSON string
	tmpStr, _ := json.Marshal(params)
	// init and format request params to string
	paramsStr = string(tmpStr[1 : len(tmpStr)-1])
	jsonStr = fmt.Sprintf("{%s, %s}", jsonStr, paramsStr)

	// step4. create Request object
	if req, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr))); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", hdler.UserAgent)
	req.Header.Set("Timeout", hdler.Timeout)

	// step5. access speech recognition API by HTTP
	if resp, err = hdler.Client.Do(req); err != nil {
		//log.Fatal(e)
		return
	}
	// step6. serialize to result string
	if result, statusCode, err = hdler.processResp(resp); err != nil {
		//log.Fatal(e)
		return
	}
	return
}

// Recognize is the major method for initiating a recognition request
func (hdler *Handler) Recognize(secretID string, dataInfoSlice []*tupumodel.DataInfo, tasks []string) (result string, statusCode int, e error) {
	// Only 10 data can be carried in one request
	if tupuerrorlib.StringIsEmpty(secretID) {
		result = ""
		statusCode = 400
		e = fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}

	var (
		url    = hdler.apiURL + secretID
		req    *http.Request
		resp   *http.Response
		params map[string]string
	)

	if params, e = hdler.GetGeneralParams(secretID); e != nil {
		statusCode = 400
		return
	}

	if req, e = hdler.request(&url, &params, dataInfoSlice, tasks); e != nil {
		//log.Fatal(e)
		return
	}

	if resp, e = hdler.Client.Do(req); e != nil {
		//log.Fatal(e)
		return
	}

	if result, statusCode, e = hdler.processResp(resp); e != nil {
		//log.Fatal(e)
		return
	}
	//fmt.Println(resp.Header)
	return
}

// GetGeneralParams is general function for getting TUPU base params
func (hdler *Handler) GetGeneralParams(secretID string) (map[string]string, error) {
	if tupuerrorlib.StringIsEmpty(secretID) {
		return nil, fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}

	var (
		signature string
		e         error
		t         = time.Now()
		timestamp = strconv.FormatInt(t.Unix(), 10)
		r         = rand.New(rand.NewSource(t.UnixNano()))
		nonce     = strconv.FormatInt(int64(r.Uint32()), 10)
		forSign   = strings.Join([]string{secretID, timestamp, nonce}, ",")
		params    = map[string]string{
			"timestamp": timestamp,
			"nonce":     nonce,
		}
	)

	if signature, e = hdler.sign([]byte(forSign)); e != nil {
		return nil, e
	}

	params["signature"] = signature

	if len(hdler.UID) > 0 {
		params["uid"] = hdler.UID
	}
	return params, nil
}

func (hdler *Handler) sign(message []byte) (string, error) {
	signed, e := hdler.signer.Sign(message)
	if e != nil {
		return "", fmt.Errorf("could not sign message: %v", e)
	}
	return base64.StdEncoding.EncodeToString(signed), nil
}

func (hdler *Handler) verify(message []byte, sig string) error {
	data, e := base64.StdEncoding.DecodeString(sig)
	if e != nil {
		return fmt.Errorf("could not decode with Base64: %v", e)
	}

	e = hdler.verifier.Verify(message, data)
	if e != nil {
		return fmt.Errorf("could not verify request: %v", e)
	}
	return nil
}

func (hdler *Handler) request(url *string, params *map[string]string, dataInfoSlice []*tupumodel.DataInfo, tasks []string) (req *http.Request, e error) {
	// verify legatity params
	if tupuerrorlib.PtrIsNil(url, params, dataInfoSlice) {
		return nil, fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}

	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)

	for key, val := range *params {
		_ = writer.WriteField(key, val)
	}

	for _, task := range tasks {
		_ = writer.WriteField("task", task)
	}

	// write binary data to request body
	for index, dataInfoItem := range dataInfoSlice {
		if e = addDataInfoField(writer, dataInfoItem, index); e == nil {
			// with other message
			if dataInfoItem.OtherMsg != nil {
				for kStr, v := range dataInfoItem.OtherMsg {
					_ = writer.WriteField(kStr, v)
				}
			}
		}
	}

	if e = writer.Close(); e != nil {
		return
	}

	// construct request object
	if req, e = http.NewRequest("POST", *url, body); e != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", hdler.UserAgent)
	req.Header.Set("Timeout", hdler.Timeout)
	// fmt.Println(req)

	return
}

func addDataInfoField(writer *multipart.Writer, dataInfo *tupumodel.DataInfo, idx int) (e error) {
	// verify legatity params
	if tupuerrorlib.PtrIsNil(writer, dataInfo) {
		return fmt.Errorf("[Params ERROR]: *io.writer or *dataInfo is null")
	}

	switch {
	case len(dataInfo.RemoteInfo) > 0:
		_ = writer.WriteField(dataInfo.FileType, dataInfo.RemoteInfo)
	case len(dataInfo.Path) > 0:
		var (
			file *os.File
			part io.Writer
		)
		if file, e = os.Open(dataInfo.Path); e != nil {
			return
		}
		part, e = writer.CreateFormFile(dataInfo.FileType, filepath.Base(dataInfo.Path))
		if e == nil {
			_, e = io.Copy(part, file)
		}
		file.Close()
	case dataInfo.Buf != nil && dataInfo.Buf.Len() > 0 && len(dataInfo.FileName) > 0:
		var part io.Writer
		part, e = writer.CreateFormFile(dataInfo.FileType, dataInfo.FileName)
		if e == nil {
			_, e = io.Copy(part, dataInfo.Buf)
		}
	default:
		return fmt.Errorf("invalid data resource at index [%v]", idx)
	}
	return
}

func (hdler *Handler) processResp(resp *http.Response) (result string, statusCode int, e error) {
	statusCode = resp.StatusCode
	//if resp.StatusCode > 500 {
	//	if resp.StatusCode == 502 {
	//		e = ErrBadGateway
	//	} else if resp.StatusCode == 503 {
	//		e = ErrServiceUnavailable
	//	} else {
	//		e = errors.New(resp.Status)
	//	}
	//	return
	//}

	body := &bytes.Buffer{}
	if _, e = body.ReadFrom(resp.Body); e != nil {
		return
	}
	if e = resp.Body.Close(); e != nil {
		return
	}

	var (
		data map[string]string
		ok   bool
		sig  string
	)
	if err := json.Unmarshal(body.Bytes(), &data); err != nil {
		if statusCode == 400 || statusCode <= 299 {
			e = errors.New("missing valid response body")
		}
		return
	} else if result, ok = data["json"]; !ok {
		e = errors.New("no result string")
		return
	} else if sig, ok = data["signature"]; !ok {
		e = errors.New("no server signature")
		return
	}
	e = hdler.verify([]byte(result), sig)
	return
}
