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
	Client *http.Client
}

// NewHandlerWithURL is also an initializer for a Handler
func NewHandlerWithURL(privateKeyPath, url string) (h *Handler, e error) {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(privateKeyPath, url) {
		return nil, fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}
	h = new(Handler)

	// init other default proprety
	h.initHandler()
	h.apiURL = RootAPIURL
	h.apiURL = url

	if h.verifier, e = tuputools.LoadTupuPublicKey(); e != nil {
		return nil, e
	}
	if h.signer, e = tuputools.LoadPrivateKey(privateKeyPath); e != nil {
		return nil, e
	}
	return h, nil
}

func (h *Handler) initHandler() {
	h.UserAgent = DefaultUserAgent
	h.ContentType = DefaultContentType
	h.Timeout = "30"
	h.Client = &http.Client{}
}

// RecognizeWithJSON is one of major method to access recognition api
func (h *Handler) RecognizeWithJSON(jsonStr, secretID, timeOut string) (result string, statusCode int, err error) {

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
		url       = h.apiURL + secretID
		req       *http.Request
		resp      *http.Response
	)

	// step2. get timestamp, nonce, signature
	if params, err = h.GetGeneralParams(secretID); err != nil {
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
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Timeout", timeOut)

	// step5. access speech recognition API by HTTP
	if resp, err = h.Client.Do(req); err != nil {
		//log.Fatal(e)
		return
	}
	// step6. serialize to result string
	if result, statusCode, err = h.processResp(resp); err != nil {
		//log.Fatal(e)
		return
	}
	return
}

// Recognize is the major method for initiating a recognition request
func (h *Handler) Recognize(secretID string, dataInfoSlice []*tupumodel.DataInfo) (result string, statusCode int, e error) {
	// Only 10 data can be carried in one request
	if len(dataInfoSlice) > 10 || tupuerrorlib.StringIsEmpty(secretID) {
		result = ""
		statusCode = 400
		e = fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName())
	}

	var (
		url    = h.apiURL + secretID
		req    *http.Request
		resp   *http.Response
		params map[string]string
	)

	if params, e = h.GetGeneralParams(secretID); e != nil {
		statusCode = 400
		return
	}

	if req, e = h.request(&url, &params, dataInfoSlice); e != nil {
		//log.Fatal(e)
		return
	}

	fmt.Println("---------------------------------------------------")
	fmt.Println(req)
	fmt.Println("---------------------------------------------------")

	if resp, e = h.Client.Do(req); e != nil {
		//log.Fatal(e)
		return
	}

	fmt.Println("---------------------------------------------------")
	fmt.Println(resp)
	fmt.Println("---------------------------------------------------")
	if result, statusCode, e = h.processResp(resp); e != nil {
		//log.Fatal(e)
		return
	}
	//fmt.Println(resp.Header)
	return
}

// GetGeneralParams is general function for getting TUPU base params
func (h *Handler) GetGeneralParams(secretID string) (map[string]string, error) {
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
			"signature": signature,
		}
	)

	if signature, e = h.sign([]byte(forSign)); e != nil {
		return nil, e
	}

	if len(h.UID) > 0 {
		params["uid"] = h.UID
	}
	return params, nil
}

func (h *Handler) sign(message []byte) (string, error) {
	signed, e := h.signer.Sign(message)
	if e != nil {
		return "", fmt.Errorf("could not sign message: %v", e)
	}
	return base64.StdEncoding.EncodeToString(signed), nil
}

func (h *Handler) verify(message []byte, sig string) error {
	data, e := base64.StdEncoding.DecodeString(sig)
	if e != nil {
		return fmt.Errorf("could not decode with Base64: %v", e)
	}

	e = h.verifier.Verify(message, data)
	if e != nil {
		return fmt.Errorf("could not verify request: %v", e)
	}
	return nil
}

func (h *Handler) request(url *string, params *map[string]string, dataInfoSlice []*tupumodel.DataInfo) (req *http.Request, e error) {
	// verify legatity params
	if tupuerrorlib.PtrIsNil(url, params, dataInfoSlice) {
		return nil, fmt.Errorf("%s, %s", tupuerrorlib.ErrorParamsIsEmpty, tupuerrorlib.GetCallerFuncName)
	}

	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)

	for key, val := range *params {
		_ = writer.WriteField(key, val)
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
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Timeout", h.Timeout)

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

func (h *Handler) processResp(resp *http.Response) (result string, statusCode int, e error) {
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
	e = h.verify([]byte(result), sig)
	return
}