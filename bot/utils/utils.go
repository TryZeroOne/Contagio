package utils

import (
	"bytes"
	rnd "crypto/rand"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var StopChan = make(chan bool)
var ClientStopChan = make(chan int)

func RandomString(strlen int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	const chars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}
	return string(result)
}

func RandomInt(strlen int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	const chars = "1234567890"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[r1.Intn(len(chars))]
	}
	conv, _ := strconv.Atoi(string(result))

	return conv
}

func GetUserAgent() string {

	uag := strings.Join(UserAgents, "\n")
	strs := strings.Split(uag, "\n")
	cStrs := len(strs)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	d := strs[r1.Intn(cStrs)]

	return d

}

func GetArrayVal(arr []string) string {

	list := strings.Join(arr, "\n")
	strs := strings.Split(list, "\n")
	cStrs := len(strs)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	d := strs[r1.Intn(cStrs)]

	return d

}

// 1 - null
// 2 - default (random bytes)
func BuildPayload(p, size int) (packet []byte) {

	rand.NewSource(time.Now().UnixNano())

	switch p {
	case 1: // null
		payload := make([]byte, size)

		randbytes := make([]byte, 30)

		_, err := rnd.Read(randbytes)
		if err != nil {
			return
		}

		payload = append(randbytes, payload...)
		payload = append(payload, randbytes...)

		packet = payload

	case 2: // default

		payload := make([]byte, size)

		_, err := rnd.Read(payload)
		if err != nil {
			return
		}

		packet = payload
	}

	return packet
}

func FormatCommand(command string) []byte {

	var res []byte

	res = bytes.ReplaceAll([]byte(command), []byte{4}, []byte{})

	return res

}

func Checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data); i += 2 {
		sum += uint32(data[i+1])<<8 | uint32(data[i])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	return ^uint16(sum)
}
