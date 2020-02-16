package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

//Encrypt doc
//Summary aes encrypt
//@Param (string) key
//@Param (string) plaintext
//@Return (string)  ciphertext
//@Return (error)
func Encrypt(key string, plaintext string) (string, error) {
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

//Decrypt doc
//@Summary aes decrypt
//@Param (string) key
//@Param (string) ciphertext
//@Return (string) plaintext
//@Return (error)
func Decrypt(key string, d string) (string, error) {
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
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
