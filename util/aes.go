package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

//AesEncrypt desc
//@Method AesEncrypt desc: aes encrypt
//@Param (string) key
//@Param (string) plaintext
//@Return (string)  ciphertext
//@Return (error)
func AesEncrypt(key string, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil
}

//AesDecrypt desc
//@Method AesDecrypt desc: aes decrypt
//@Param (string) key
//@Param (string) ciphertext
//@Return (string) plaintext
//@Return (error)
func AesDecrypt(key string, d string) (string, error) {
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	fmt.Println(len(ciphertext), len(iv))
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
