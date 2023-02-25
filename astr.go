package security

import (
	"math/rand"
	"time"
	"unsafe"
)

func GenerateRandStr(n int) string {
	var lettersLen int = 69
	var letters = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&()")
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]byte, n)
	for i := range r {
		r[i] = letters[rd.Intn(lettersLen)]
	}
	return Bytes2Str(r)
}

func Str2SliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
