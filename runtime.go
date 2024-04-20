package easyuse

import "runtime"

// GetCurrentGoroutineStack 获取当前写成跟踪信息
func GetCurrentGoroutineStack(bufByte int64) string {
	var buf = make([]byte, 0, bufByte)
	n := runtime.Stack(buf, false)
	return Bytes2Str(buf[:n])
}
