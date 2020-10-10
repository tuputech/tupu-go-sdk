package errorlib

import (
	"reflect"
	"runtime"
)

const (
	// ErrorParamsIsEmpty describes illegal parameters
	ErrorParamsIsEmpty = "Params can't empty"
)

// GetCurrentFuncName is a util function for gettiing current execution function name
func GetCurrentFuncName() string {
	return getFuncName(1)
}

// GetCallerFuncName is a util function for getting caller function name
func GetCallerFuncName() string {
	return getFuncName(2)
}

func getFuncName(rank int) string {
	pc, _, _, _ := runtime.Caller(rank)
	return runtime.FuncForPC(pc).Name()
}

// PtrIsNil is a util function for check the legality of the params
func PtrIsNil(data ...interface{}) bool {
	for _, v := range data {
		vPtr := reflect.ValueOf(v)
		if (vPtr.Kind() == reflect.Ptr) && vPtr.IsNil() {
			return true
		}
	}
	return false
}

// StringIsEmpty is a util function for check the empty of the string
func StringIsEmpty(str ...string) bool {
	for _, v := range str {
		if len(v) == 0 {
			return true
		}
	}
	return false
}
