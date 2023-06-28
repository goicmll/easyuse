package easyuse

import "runtime"


func GetCurrentGoroutineStack(bufByte int64) string {
	var buf = make([]byte, 0, bufByte)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
