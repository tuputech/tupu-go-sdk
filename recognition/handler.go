package recognition

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	tupucontrol "github.com/tuputech/tupu-go-sdk/lib/controller"
	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

var (
	ErrBadGateway         = errors.New("502 Bad Gateway")
	ErrServiceUnavailable = errors.New("503 Service Unavailable")
	ImageRecognitionURL   = "http://api.open.tuputech.com/v3/recognition/"
)

type (
	config struct {
		tags  []string
		tasks []string
	}
	options func(*config)

	// Handler is a client-side helper to access TUPU visual recognition service
	Handler struct {
		hdler   *tupucontrol.Handler
		Client  *http.Client
		imgPool sync.Pool
		//
		UID       string //for sub-user statistics and billing
		UserAgent string
	}
)

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

func (h *Handler) WithTags(tags []string) options {
	return func(c *config) {
		c.tags = tags
	}
}

func (h *Handler) WithTasks(tasks []string) options {
	return func(c *config) {
		c.tasks = tasks
	}
}

// PerformWithURL is a shortcut for initiating a recognition request with URLs of images
func (h *Handler) PerformWithURL(secretID string, imageURLs []string, options ...func(*config)) (result string, statusCode int, e error) {
	var images []*Image
	for index, val := range imageURLs {
		img := h.imgPool.Get().(*Image)
		img.InitConf(tupumodel.WithFileURL(val))
		images[index] = img
		defer h.recycleDataObj(img)
	}
	var c config
	for _, fn := range options {
		fn(&c)
	}
	return h.Perform(secretID, images, c.tags, c.tasks)
}

// PerformWithPath is a shortcut for initiating a recognition request with paths of images
func (h *Handler) PerformWithPath(secretID string, imagePaths []string, options ...func(*config)) (result string, statusCode int, e error) {
	images := make([]*Image, len(imagePaths))
	for i, val := range imagePaths {
		img := h.imgPool.Get().(*Image)
		img.InitConf(tupumodel.WithLocalPath(val))
		images[i] = img
		defer h.recycleDataObj(img)
	}
	var c config
	for _, fn := range options {
		fn(&c)
	}
	return h.Perform(secretID, images, c.tags, c.tasks)
}

// Perform is the major method for initiating a recognition request
func (h *Handler) Perform(secretID string, images []*Image, tags []string, tasks []string) (result string, statusCode int, e error) {
	// verify legatity params
	if tupuerror.PtrIsNil(images) || tupuerror.StringIsEmpty(secretID) {
		e = fmt.Errorf("%s, %s", tupuerror.ErrorParamsIsEmpty, tupuerror.GetCallerFuncName())
	}

	var (
		tagsLen       = len(tags)
		imagesLen     = len(images)
		dataInfoSlice = make([]*tupumodel.DataInfo, imagesLen+1)
	)

	for i := 0; i < imagesLen; i++ {
		images[i].dataInfo.OtherMsg = make(map[string]string)
		if i < tagsLen {
			images[i].dataInfo.OtherMsg["tag"] = tags[i]
		}
		dataInfoSlice[i] = images[i].dataInfo
	}

	return h.hdler.Recognize(secretID, dataInfoSlice, tasks)

}

func (h *Handler) recycleDataObj(img *Image) {
	img.dataInfo.ClearData()
	h.imgPool.Put(img)
}
