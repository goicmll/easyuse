package habits

import (
	"math/rand"
	"time"
)

var L1 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&()")
var L2 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var L3 = Str2SliceByte("0123456789")

func GenerateRandStr(n int, source []byte) string {
	var lettersLen = len(source)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]byte, n)
	for i := range r {
		r[i] = source[rd.Intn(lettersLen)]
	}
	return Bytes2Str(r)
}
