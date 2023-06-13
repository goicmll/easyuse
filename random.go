package habits

import (
	"math/rand"
	"time"
)

// 预定义随机字符串源
var S1 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#%^&()")
var S2 = Str2SliceByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var S3 = Str2SliceByte("abcdefghijklmnopqrstuvwxyz")
var S4 = Str2SliceByte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var S5 = Str2SliceByte("0123456789")

// 生成随机字符传
func Random(n int, source []byte) string {
	var lettersLen = len(source)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := make([]byte, n)
	for i := range r {
		r[i] = source[rd.Intn(lettersLen)]
	}
	return Bytes2Str(r)
}
