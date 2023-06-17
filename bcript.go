package eu

import "golang.org/x/crypto/bcrypt"

// BcryptStr 字符串生成bcrypt
func BcryptStr(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(Str2SliceByte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", NewHabitError(err.Error())
	}
	return Bytes2Str(b), nil
}

// BcryptCompare 比较密文和字符串
func BcryptCompare(strHash, str string) bool {
	byteHash := Str2SliceByte(strHash)
	bytePwd := Str2SliceByte(str)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
