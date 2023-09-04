package easyuse

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

// AesCBCStrEncrypt key长度必须为16, 24或者32, 返回hex.EncodeToString 后的字符串
// 这边IV直接就使用 key
func AesCBCStrEncrypt(key, origin string) (string, error) {
	// 转成字节数组
	originByte := Str2SliceByte(origin)
	keyByte := Str2SliceByte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", NewEasyUseError("key 不可用, 必须为16, 24, 32长度的字符!")
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	pk := pkcs7Padding(originByte, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	// 创建数组
	out := make([]byte, len(pk))
	// 加密
	blockMode.CryptBlocks(out, pk)
	return hex.EncodeToString(out), nil
}

// key长度必须为16, 24或者32, 参数origin 为hex.EncodeToString 后的字符串
// 这边IV直接就使用 key
func AesCBCStrDecrypt(key, origin string) (string, error) {
	// 转成字节数组
	originByte, err := hex.DecodeString(origin)
	if err != nil {
		return "", NewEasyUseError("密文错误, 无法解密!")
	}
	keyByte := Str2SliceByte(key)
	// 分组秘钥
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", NewEasyUseError("key 不可用, 必须为16, 24, 32长度的字符!")
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	// 创建数组
	out := make([]byte, len(originByte))
	// 解密
	blockMode.CryptBlocks(out, originByte)
	//去除加密是补全的字节, 加密是把补全的位数存入了加密字符串字节切片的最后几位
	// 去补全码
	out, err = pkcs7UnPadding(out)
	return Bytes2Str(out), err
}

// ########### AES CRT
// key长度必须为16, 24或者32
func AesCRTCrypt(key, origin string) (string, error) {
	// 转成字节数组
	plainText := Str2SliceByte(origin)
	keyByte := Str2SliceByte(key)
	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", NewEasyUseError("key 不可用, 必须为16, 24, 32长度的字符!")
	}
	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat(Str2SliceByte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return Bytes2Str(dst), nil
}

// ########### AES OFB
// key长度必须为16, 24或者32
func AesOFBStrEncrypt(key, origin string) (string, error) {
	// 转成字节数组
	originByte := Str2SliceByte(origin)
	keyByte := Str2SliceByte(key)
	originByte = pkcs7Padding(originByte, aes.BlockSize)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", NewEasyUseError("key 不可用, 必须为16, 24, 32长度的字符!")
	}
	out := make([]byte, aes.BlockSize+len(originByte))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], originByte)
	return Bytes2Str(out), nil
}

// key长度必须为16, 24或者32
func AesOFBStrDecrypt(key, origin string) (string, error) {
	originByte := Str2SliceByte(origin)
	keyByte := Str2SliceByte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", NewEasyUseError("key 不可用, 必须为16, 24, 32长度的字符!")
	}
	iv := originByte[:aes.BlockSize]
	originByte = originByte[aes.BlockSize:]
	if len(originByte)%aes.BlockSize != 0 {
		return "", NewEasyUseError("解密失败: data is not a multiple of the block size")
	}

	out := make([]byte, len(originByte))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(out, originByte)

	out, err = pkcs7UnPadding(out)
	return Bytes2Str(out), err
}

// 补码
// AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
// 补全blockSize的倍数个数, 即len(pk)%blockSize=0
func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去码
func pkcs7UnPadding(origin []byte) ([]byte, error) {
	length := len(origin)
	if length < 1 {
		return Str2SliceByte(""), NewEasyUseError("解密失败!")
	}
	unpadding := int(origin[length-1])
	//这边取最后一位就一定是加密时补全的位数
	if length < unpadding {
		return Str2SliceByte(""), NewEasyUseError("解密失败!")
	}
	return origin[:(length - unpadding)], nil
}

// BcryptStr 字符串生成bcrypt
func BcryptStr(str string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(Str2SliceByte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", NewEasyUseError(err.Error())
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

// JwtMake 制作jwt串
func JwtMake(id, secret, issuer, subject, audience string, maxAge int) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  []string{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(maxAge) * time.Second)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtStr, err := token.SignedString(Str2SliceByte(secret))
	if err != nil {
		return "", NewEasyUseError(err.Error())
	}
	return jwtStr, nil
}

// JwtParse 解析jwt
func JwtParse(jwtStr, secret string) (*jwt.RegisteredClaims, error) {
	tok, err := jwt.ParseWithClaims(jwtStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return Str2SliceByte(secret), nil
	})
	if err != nil {
		return nil, NewEasyUseError("token 不可用!")
	}
	if claims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, NewEasyUseError("token 不可用!")
}

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

// StringMd5 md5
func StringMd5(s string) string {
	data := Str2SliceByte(s) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
