package model

import (
	"bytes"

	tupuerrorlib "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

// DataInfo is a wrapper for storing url, path or binary content of the data
type DataInfo struct {
	// FileType Identifies this file type for request parameters
	FileType string
	// RemoteInfo is the remote address of the file
	RemoteInfo string
	// Path is the local path of the file
	Path string
	// FileName is the rename of the file by caller user
	FileName string
	// Buf is the binary data of the file
	Buf *bytes.Buffer
	// OtherMsg is the other accompanying message
	OtherMsg map[string]string
}

// NewRemoteDataInfo is an initializer for create data resource with a url
func NewRemoteDataInfo(remoteInfo string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(remoteInfo) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.RemoteInfo = remoteInfo
	return dataInfo
}

// NewLocalDataInfo is an initializer for create data resource with a file path
func NewLocalDataInfo(path string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(path) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.Path = path
	return dataInfo
}

// NewBinaryDataInfo is an initializer for create data resource with binary content
func NewBinaryDataInfo(buf []byte, filename string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(filename) || tupuerrorlib.PtrIsNil(buf) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.Buf = bytes.NewBuffer(buf)
	dataInfo.FileName = filename
	return dataInfo
}

// Tag is an helper to set property tag of DataInfo and return DataInfo itself
// func (i *DataInfo) Tag(tag string) *DataInfo {
// 	i.tag = tag
// 	return i
// }

// ClearBuffer is an helper to set property tag of DataInfo and return DataInfo itself
func (dataInfo *DataInfo) ClearBuffer() {
	dataInfo.Buf = nil
}
