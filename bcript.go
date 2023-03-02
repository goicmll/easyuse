package security

import "golang.org/x/crypto/bcrypt"

// 字符串生成bcrity
func BcryptStr(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(Str2SliceByte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", NewSecurityError(err.Error())
	}
	return Bytes2Str(b), nil
}

// 比较密文和字符串
func BcryptCompare(strHash, str string) bool {
	byteHash := Str2SliceByte(strHash)
	bytePwd := Str2SliceByte(str)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
