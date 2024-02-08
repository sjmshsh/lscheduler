package util

import (
	"context"
	"encoding/json"
	"runtime"
)

const (
	panicStackSize = 1 << 16
)

// GetMarshalStr .
func GetMarshalStr(obj interface{}) (result string) {
	raw, _ := json.Marshal(obj)
	result = string(raw)
	return
}

// PanicRecover ...
func PanicRecover(ctx context.Context, recoverResult interface{}) bool {
	if recoverResult == nil {
		return false
	}

	buf := make([]byte, panicStackSize)
	buf = buf[:runtime.Stack(buf, false)]
	return true
}
