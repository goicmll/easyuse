package easyuse

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// RsaDecrypt 解密
func RsaDecrypt(originBase64, privateKey string) (string, error) {
	originByte, err := base64.StdEncoding.DecodeString(originBase64)
	if err != nil {
		return "", NewEasyUseError("密文base64解码错误！")
	}
	keyByte := []byte(privateKey)
	//解密
	block, _ := pem.Decode(keyByte)
	if block == nil {
		return "", NewEasyUseError("RSA 私钥解码错误！")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", NewEasyUseError("RSA 解析私钥错误！")
	}
	// 解密
	out, err := rsa.DecryptPKCS1v15(rand.Reader, priv, originByte)
	if err != nil {
		return "", NewEasyUseError("RSA 密文解密错误")
	}
	return Bytes2Str(out), nil
}

// RsaEncrypt ########### RSA 加解密
func RsaEncrypt(origin, publicKey string) (string, error) {
	originByte := []byte(origin)
	keyByte := []byte(publicKey)
	//解密pem格式的公钥
	block, _ := pem.Decode(keyByte)
	if block == nil {
		return "", NewEasyUseError("RSA 公钥解码错误！")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", NewEasyUseError("RSA 解析公钥错误！")
	}
	// 类型断言
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return "", NewEasyUseError("RSA 解析公钥错误！")
	}
	//加密
	out, err := rsa.EncryptPKCS1v15(rand.Reader, pub, originByte)
	if err != nil {
		return "", NewEasyUseError("RSA 密文加密错误！")
	}
	return base64.StdEncoding.EncodeToString(out), nil
}
