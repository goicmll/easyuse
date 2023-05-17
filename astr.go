package security

import (
	"math/rand"
	"time"
	"unsafe"
)

func GenerateRandStr(n int) string {
	var lettersLen = 69
	var letters = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&()")
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]byte, n)
	for i := range r {
		r[i] = letters[rd.Intn(lettersLen)]
	}
	return Bytes2Str(r)
}

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
