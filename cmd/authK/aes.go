package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AesCryptor struct {
	key []byte
	iv  []byte
}

func New(akey []byte) *AesCryptor {
	iv := make([]byte, aes.BlockSize)
	key := make([]byte, aes.BlockSize)
	copy(key, akey)
	return &AesCryptor{key, iv}
}

// 加密数据
func (a *AesCryptor) Encrypt(data []byte) ([]byte, error) {
	aesBlockEncrypter, err := aes.NewCipher(a.key)
	content := pkcs5Padding(data, aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))
	if err != nil {
		println(err.Error())
		return nil, err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, a.iv)
	aesEncrypter.CryptBlocks(encrypted, content)
	return encrypted, nil
}

// 解密数据
func (a *AesCryptor) Decrypt(src []byte) (data []byte, err error) {
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(a.key)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, a.iv)
	aesDecrypter.CryptBlocks(decrypted, src)
	return pkcs5Trimming(decrypted), nil
}

// PKCS5包装，如果刚好是16的倍数，就多填充一个block
func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// 解包装
func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
