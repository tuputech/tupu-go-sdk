package recognition

import (
	"errors"
	"fmt"
	"net/http"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

var (
	ErrBadGateway         = errors.New("502 Bad Gateway")
	ErrServiceUnavailable = errors.New("503 Service Unavailable")
	ImageRecognitionURL   = "http://api.open.tuputech.com/v3/recognition/"
)

// Handler is a client-side helper to access TUPU visual recognition service
type Handler struct {
	hdler *tupucontrol.Handler
	//
	UID       string //for sub-user statistics and billing
	UserAgent string
	Client    *http.Client
}

// NewHandler is an initializer for a Handler
func NewHandler(privateKeyPath string) (*Handler, error) {
	var (
		e error
		h = new(Handler)
	)
	h.hdler, e = tupucontrol.NewHandlerWithURL(privateKeyPath, ImageRecognitionURL)
	if e != nil {
		return nil, e
	}
	return h, nil
}

// NewHandlerWithURL is also an initializer for a Handler
func NewHandlerWithURL(privateKeyPath, url string) (h *Handler, e error) {
	h = new(Handler)
	h.hdler, e = tupucontrol.NewHandlerWithURL(privateKeyPath, url)
	if e != nil {
		return nil, e
	}
	return h, nil
}

// PerformWithURL is a shortcut for initiating a recognition request with URLs of images
func (h *Handler) PerformWithURL(secretID string, imageURLs []string, tags []string, tasks []string) (result string, statusCode int, e error) {
	var images []*Image
	for _, val := range imageURLs {
		images = append(images, NewRemoteImage(val))
	}
	return h.Perform(secretID, images, tags, tasks)
}

// PerformWithPath is a shortcut for initiating a recognition request with paths of images
func (h *Handler) PerformWithPath(secretID string, imagePaths []string, tags []string, tasks []string) (result string, statusCode int, e error) {
	var images []*Image
	for _, val := range imagePaths {
		images = append(images, NewLocalImage(val))
	}
	return h.Perform(secretID, images, tags, tasks)
}

// Perform is the major method for initiating a recognition request
func (h *Handler) Perform(secretID string, images []*Image, tags []string, tasks []string) (result string, statusCode int, e error) {
	// verify legatity params
	if tupuerror.PtrIsNil(images) || tupuerror.StringIsEmpty(secretID) {
		e = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	// once request only can carry in 10 images
	if len(images) > 10 {
		e = fmt.Errorf("Once request only can bring 10 images")
		return
	}

	var (
		tagsLen       = len(tags)
		imagesLen     = len(images)
		dataInfoSlice = make([]*tupumodel.DataInfo, 0)
	)

	for i := 0; i < imagesLen; i++ {
		images[i].dataInfo.OtherMsg = make(map[string]string)
		if i < tagsLen {
			images[i].dataInfo.OtherMsg["tag"] = tags[i]
		}
		dataInfoSlice = append(dataInfoSlice, images[i].dataInfo)
	}

	return h.hdler.Recognize(secretID, dataInfoSlice, tasks)

}
