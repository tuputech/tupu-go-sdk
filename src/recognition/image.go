package recognition

import "bytes"

// Image is a wrapper for storing url, path or binary content of an image
type Image struct {
	url      string
	path     string
	buf      *bytes.Buffer
	filename string
	tag      string
}

// NewRemoteImage is an initializer for create image resource with a url
func NewRemoteImage(url string) *Image {
	i := new(Image)
	i.url = url
	return i
}

// NewLocalImage is an initializer for create image resource with a file path
func NewLocalImage(path string) *Image {
	i := new(Image)
	i.path = path
	return i
}

// NewBinaryImage is an initializer for create image resource with binary content
func NewBinaryImage(buf []byte, filename string) *Image {
	i := new(Image)
	i.buf = bytes.NewBuffer(buf)
	i.filename = filename
	return i
}

// Tag is an helper to set property tag of Image and return Image itself
func (i *Image) Tag(tag string) *Image {
	i.tag = tag
	return i
}

// ClearBuffer is an helper to set property tag of Image and return Image itself
func (i *Image) ClearBuffer() {
	i.buf = nil
}
