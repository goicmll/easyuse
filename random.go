package habits

import (
	"math/rand"
	"time"
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
