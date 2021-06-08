package recognition

import (
	"bytes"

	tupuerror "github.com/tuputech/tupu-go-sdk/lib/errorlib"
	tupumodel "github.com/tuputech/tupu-go-sdk/lib/model"
)

// Image is a wrapper for storing url, path or binary content of an image
type (
	Image struct {
		dataInfo *tupumodel.DataInfo
	}

	imgOptFunc func(*Image)
)

func newImage() *Image {
	img := new(Image)
	img.dataInfo = new(tupumodel.DataInfo)
	img.dataInfo.SetFileType("image")
	return img
}

// NewRemoteImage is an initializer for create image resource with a url
func NewRemoteImage(url string) *Image {
	// verify legatity params
	if tupuerror.StringIsEmpty(url) {
		return nil
	}

	img := newImage()
	img.dataInfo.SetRemoteInfo(url)
	return img
}

// NewLocalImage is an initializer for create image resource with a file path
func NewLocalImage(path string) *Image {
	// verfify legatity params
	if tupuerror.StringIsEmpty(path) {
		return nil
	}

	img := newImage()
	img.dataInfo.SetPath(path)
	return img
}

// NewBinaryImage is an initializer for create image resource with binary content
func NewBinaryImage(buf []byte, filename string) *Image {
	// verify legatity params
	if tupuerror.StringIsEmpty(filename) || tupuerror.PtrIsNil(buf) {
		return nil
	}

	img := newImage()
	img.dataInfo.Buf = bytes.NewBuffer(buf)
	img.dataInfo.FileName = filename
	return img
}

// Tag is an helper to set property tag of Image and return Image itself
// func (i *Image) Tag(tag string) *Image {
// 	i.tag = tag
// 	return i
// }

// ClearBuffer is an helper to set property tag of Image and return Image itself
func (img *Image) ClearBuffer() {
	img.dataInfo.Buf = nil
}

func (img *Image) InitConf(optFuncs ...tupumodel.OptFunc) {
	for _, optFunc := range optFuncs {
		optFunc(img.dataInfo)
	}
}
