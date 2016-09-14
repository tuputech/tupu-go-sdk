package recognition

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Handler is a client-side helper to access TUPU visual recognition service
type Handler struct {
	apiURL   string
	signer   Signer
	verifier Verifier
	//
	UID       string //for sub-user statistics and billing
	UserAgent string
}

// NewHandler is an initializer for a Handler
func NewHandler(privateKeyPath string) (*Handler, error) {
	h := new(Handler)
	h.apiURL = "http://api.open.tuputech.com/v3/recognition/"
	h.UserAgent = "tupu-client/1.0"

	var e error
	if h.verifier, e = LoadTupuPublicKey(); e != nil {
		return nil, e
	}
	if h.signer, e = LoadPrivateKey(privateKeyPath); e != nil {
		return nil, e
	}
	return h, nil
}

// NewHandlerWithURL is also an initializer for a Handler
func NewHandlerWithURL(privateKeyPath, url string) (h *Handler, e error) {
	if h, e = NewHandler(privateKeyPath); e != nil {
		return
	}
	h.apiURL = url
	return h, nil
}

// Perform is the major method for initiating a recognition request
func (h *Handler) Perform(secretID string, images []*Image, tags []string) (result string, statusCode int, e error) {
	t := time.Now()
	timestamp := strconv.FormatInt(t.Unix(), 10)
	r := rand.New(rand.NewSource(t.UnixNano()))
	nonce := strconv.FormatInt(int64(r.Uint32()), 10)
	forSign := strings.Join([]string{secretID, timestamp, nonce}, ",")
	var signature string
	if signature, e = h.sign([]byte(forSign)); e != nil {
		return
	}

	params := map[string]string{
		"timestamp": timestamp,
		"nonce":     nonce,
		"signature": signature,
	}
	if len(h.UID) > 0 {
		params["uid"] = h.UID
	}

	var (
		url  = h.apiURL + secretID
		req  *http.Request
		resp *http.Response
	)
	if req, e = h.request(&url, &params, images, tags); e != nil {
		log.Fatal(e)
		return
	}
	client := &http.Client{}
	if resp, e = client.Do(req); e != nil {
		log.Fatal(e)
		return
	}
	if result, e = h.processResp(resp); e != nil {
		log.Fatal(e)
		return
	}
	statusCode = resp.StatusCode
	//fmt.Println(resp.Header)
	return
}

func (h *Handler) sign(message []byte) (string, error) {
	signed, e := h.signer.Sign(message)
	if e != nil {
		return "", fmt.Errorf("Could not sign message: %v", e)
	}
	return base64.StdEncoding.EncodeToString(signed), nil
}

func (h *Handler) verify(message []byte, sig string) error {
	data, e := base64.StdEncoding.DecodeString(sig)
	if e != nil {
		return fmt.Errorf("Could not decode with Base64: %v", e)
	}

	e = h.verifier.Verify(message, data)
	if e != nil {
		return fmt.Errorf("Could not verify request: %v", e)
	}
	return nil
}

func (h *Handler) request(url *string, params *map[string]string, images []*Image, tags []string) (req *http.Request, e error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range *params {
		_ = writer.WriteField(key, val)
	}

	tagsCnt := 0
	if tags != nil {
		tagsCnt = len(tags)
	}
	tag := ""
	for i, img := range images {
		if e = addImageField(writer, img, i); e == nil {
			if i < tagsCnt {
				tag = tags[i]
			}
			if len(tag) > 0 {
				_ = writer.WriteField("tag", tag)
			}
		}
	}

	if e = writer.Close(); e != nil {
		return
	}

	req, e = http.NewRequest("POST", *url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Timeout", "30")

	fmt.Println(req.Header)
	fmt.Println(body.String())
	return
}

func addImageField(writer *multipart.Writer, img *Image, idx int) (e error) {
	switch {
	case len(img.url) > 0:
		_ = writer.WriteField("image", img.url)
	case len(img.path) > 0:
		var (
			file *os.File
			part io.Writer
		)
		if file, e = os.Open(img.path); e != nil {
			return
		}
		part, e = writer.CreateFormFile("image", filepath.Base(img.path))
		if e == nil {
			_, e = io.Copy(part, file)
		}
		file.Close()
	case img.buf != nil && img.buf.Len() > 0 && len(img.filename) > 0:
		var part io.Writer
		part, e = writer.CreateFormFile("image", img.filename)
		if e == nil {
			_, e = io.Copy(part, img.buf)
		}
	default:
		return fmt.Errorf("Invalid image resource at index [%v]", idx)
	}
	return
}

func (h *Handler) processResp(resp *http.Response) (result string, e error) {
	body := &bytes.Buffer{}
	if _, e = body.ReadFrom(resp.Body); e != nil {
		return
	}
	resp.Body.Close()

	var (
		data map[string]string
		ok   bool
		sig  string
	)
	if e = json.Unmarshal(body.Bytes(), &data); e != nil {
		return
	} else if result, ok = data["json"]; !ok {
		return "", errors.New("No result string")
	} else if sig, ok = data["signature"]; !ok {
		return "", errors.New("No server signature")
	}
	return result, h.verify([]byte(result), sig)
}
