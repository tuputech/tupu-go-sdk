package model

import (
	"bytes"

	tupuerrorlib "github.com/tuputech/tupu-go-sdk/lib/errorlib"
)

// DataInfo is a wrapper for storing url, path or binary content of the data
type DataInfo struct {
	// Buf is the binary data of the file
	Buf *bytes.Buffer
	// FileType Identifies this file type for request parameters
	FileType string
	// RemoteInfo is the remote address of the file
	RemoteInfo string
	// Path is the local path of the file
	Path string
	// FileName is the rename of the file by caller user
	FileName string
	// OtherMsg is the other accompanying message
	OtherMsg map[string]string
}

type OptFunc func(*DataInfo)

// SetOtherMsg is setting function for DataInfo object
func (dataInfo *DataInfo) SetOtherMsg(msg map[string]string) {
	if tupuerrorlib.PtrIsNil(msg) {
		return
	}
	dataInfo.OtherMsg = msg
}

// SetBuf is setting function for DataInfo object
func (dataInfo *DataInfo) SetBuf(buf []byte) {
	if tupuerrorlib.PtrIsNil(buf) {
		return
	}
	dataInfo.Buf = bytes.NewBuffer(buf)
}

// SetFileType is setting function for DataInfo object
func (dataInfo *DataInfo) SetFileType(ftype string) {
	if tupuerrorlib.StringIsEmpty(ftype) {
		return
	}
	dataInfo.FileType = ftype
}

// SetFileName is setting function for DataInfo object
func (dataInfo *DataInfo) SetFileName(fName string) {
	if tupuerrorlib.StringIsEmpty(fName) {
		return
	}
	dataInfo.FileName = fName
}

// SetRemoteInfo is setting function for DataInfo object
func (dataInfo *DataInfo) SetRemoteInfo(fRemoteInfo string) {
	if tupuerrorlib.StringIsEmpty(fRemoteInfo) {
		return
	}
	dataInfo.RemoteInfo = fRemoteInfo
}

// SetPath is setting function for DataInfo object
func (dataInfo *DataInfo) SetPath(localPath string) {
	if tupuerrorlib.StringIsEmpty(localPath) {
		return
	}
	dataInfo.Path = localPath
}

// NewRemoteDataInfo is an initializer for create data resource with a url
func NewRemoteDataInfo(remoteInfo string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(remoteInfo) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.SetRemoteInfo(remoteInfo)
	return dataInfo
}

// NewLocalDataInfo is an initializer for create data resource with a file path
func NewLocalDataInfo(path string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(path) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.SetPath(path)
	return dataInfo
}

// NewBinaryDataInfo is an initializer for create data resource with binary content
func NewBinaryDataInfo(buf []byte, filename string) *DataInfo {
	// verify legatity params
	if tupuerrorlib.StringIsEmpty(filename) || tupuerrorlib.PtrIsNil(buf) {
		return nil
	}
	dataInfo := new(DataInfo)
	dataInfo.SetBuf(buf)
	dataInfo.SetFileName(filename)
	return dataInfo
}

// ClearBuffer is an helper to set property tag of DataInfo and return DataInfo itself
func (dataInfo *DataInfo) ClearBuffer() {
	dataInfo.Buf = nil
}

// ClearData is an helper to reset DataInfo struct
func (dtInfo *DataInfo) ClearData() {
	dtInfo.Buf = nil
	dtInfo.OtherMsg = nil
	dtInfo.Path = ""
	dtInfo.FileName = ""
	dtInfo.RemoteInfo = ""
}

func WithBinary(buf []byte, filename string) OptFunc {
	return func(di *DataInfo) {
		di.Buf = bytes.NewBuffer(buf)
		di.FileName = filename
	}
}

func WithLocalPath(path string) OptFunc {
	return func(di *DataInfo) {
		di.Path = path
	}
}

func WithFileURL(uri string) OptFunc {
	return func(di *DataInfo) {
		di.RemoteInfo = uri
	}
}

func WithFileType(fileType string) OptFunc {
	return func(di *DataInfo) {
		di.FileType = fileType
	}
}
