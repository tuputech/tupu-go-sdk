package generalapi

import "bytes"

// DataInfo is a wrapper for storing url, path or binary content of the data
type DataInfo struct {
	url      string
	path     string
	filename string
	buf      *bytes.Buffer
}

// NewRemoteDataInfo is an initializer for create data resource with a url
func NewRemoteDataInfo(url string) *DataInfo {
	// verify legatity params
	if StringIsEmpty(url) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.url = url
	return dataInfo
}

// NewLocalDataInfo is an initializer for create data resource with a file path
func NewLocalDataInfo(path string) *DataInfo {
	// verify legatity params
	if StringIsEmpty(path) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.path = path
	return dataInfo
}

// NewBinaryDataInfo is an initializer for create data resource with binary content
func NewBinaryDataInfo(buf []byte, filename string) *DataInfo {
	// verify legatity params
	if StringIsEmpty(filename) || PtrIsNil(buf) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.buf = bytes.NewBuffer(buf)
	dataInfo.filename = filename
	return dataInfo
}

// Tag is an helper to set property tag of DataInfo and return DataInfo itself
// func (i *DataInfo) Tag(tag string) *DataInfo {
// 	i.tag = tag
// 	return i
// }

// ClearBuffer is an helper to set property tag of DataInfo and return DataInfo itself
func (dataInfo *DataInfo) ClearBuffer() {
	dataInfo.buf = nil
}
