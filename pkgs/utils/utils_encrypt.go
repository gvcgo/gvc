package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
)

type AesCrypt struct {
	Key []byte
}

func NewCrypt(password string) (c *AesCrypt) {
	has := md5.Sum([]byte(password))
	c = &AesCrypt{
		Key: []byte(fmt.Sprintf("%x", has)),
	}
	return
}

func (that *AesCrypt) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (that *AesCrypt) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (that *AesCrypt) AesEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(that.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = that.pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, that.Key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (that *AesCrypt) AesDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(that.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, that.Key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = that.pKCS7UnPadding(origData)
	return origData, nil
}
