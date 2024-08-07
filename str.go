package easyuse

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// S5 预定义随机字符串源
var S5 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&()")
var S4 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var S3 = Str2SliceByte("abcdefghijklmnopqrstuvwxyz")
var S2 = Str2SliceByte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var S1 = Str2SliceByte("0123456789")

// Random 生成随机字符传
func Random(n uint, source []byte) string {
	if n == 0 {
		return ""
	}
	var lettersLen = len(source)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]byte, n)
	for i := range r {
		r[i] = source[rd.Intn(lettersLen)]
	}
	return Bytes2Str(r)
}

// Str2SliceByte string to []byte
func Str2SliceByte(s string) []byte {
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// Bytes2Str bytes to string
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Placeholder 生成占位符
//
//	Placeholder(5, "?")  return: "?,?,?,?,?"
func Placeholder(n int, holder, sep string) string {
	ph := make([]string, n)
	for i := 0; i < n; i++ {
		ph[i] = holder
	}
	return strings.Join(ph, sep)
}
