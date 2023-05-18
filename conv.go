package habits

import (
	"reflect"
	"unsafe"
)

func Str2SliceByte(s string) []byte {
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
