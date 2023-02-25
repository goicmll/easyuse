package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// key长度必须为16, 24或者32, 返回hex.EncodeToString 后的字符串
// 这边IV直接就使用 key
func AesCBCStrEncrypt(key, origin string) string {
	// 转成字节数组
	originByte := []byte(origin)
	keyByte := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(keyByte)
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
	return hex.EncodeToString(out)
}

// key长度必须为16, 24或者32, 参数origin 为hex.EncodeToString 后的字符串
// 这边IV直接就使用 key
func AesCBCStrDecrypt(key, origin string) (string, error) {
	// 转成字节数组
	originByte, _ := hex.DecodeString(origin)
	keyByte := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(keyByte)
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
	out, err := pkcs7UnPadding(out)
	return string(out), err
}

// ########### AES CRT
// key长度必须为16, 24或者32
func AesCRTCrypt(key, origin string) (string, error) {
	// 转成字节数组
	plainText := []byte(origin)
	keyByte := []byte(key)
	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return string(dst), nil
}

// ########### AES OFB
// key长度必须为16, 24或者32
func AesOFBStrEncrypt(key, origin string) (string, error) {
	// 转成字节数组
	originByte := []byte(origin)
	keyByte := []byte(key)
	originByte = pkcs7Padding(originByte, aes.BlockSize)

	block, _ := aes.NewCipher(keyByte)
	out := make([]byte, aes.BlockSize+len(originByte))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], originByte)
	return string(out), nil
}

// key长度必须为16, 24或者32
func AesOFBStrDecrypt(key, origin string) (string, error) {
	originByte := []byte(origin)
	keyByte := []byte(key)
	block, _ := aes.NewCipher(keyByte)
	iv := originByte[:aes.BlockSize]
	originByte = originByte[aes.BlockSize:]
	if len(originByte)%aes.BlockSize != 0 {
		return "", SecurityError{Msg: "解密失败: data is not a multiple of the block size"}
	}

	out := make([]byte, len(originByte))
	mode := cipher.NewOFB(block, iv)
	mode.XORKeyStream(out, originByte)

	out, err := pkcs7UnPadding(out)
	return string(out), err
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
		return []byte(""), SecurityError{Msg: "解密失败!"}
	}
	unpadding := int(origin[length-1])
	//这边取最后一位就一定是加密时补全的位数
	if length < unpadding {
		return []byte(""), SecurityError{Msg: "解密失败!"}
	}
	return origin[:(length - unpadding)], nil
}
