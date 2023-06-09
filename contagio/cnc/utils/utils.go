package utils

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

func Sha3(str string) string {
	h := sha3.New512()
	h.Write([]byte(str))
	hash := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash[:])

}

func Reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return result
}
