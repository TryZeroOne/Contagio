package cnc

import (
	"bufio"
	"bytes"
	"contagio/contagio/cnc/database"
	"contagio/contagio/cnc/utils"
	"fmt"
	random "math/rand"
	"net/textproto"
	"strings"
	"time"
)

func (c *Connection) Auth() bool {

	c.conn.SetDeadline(time.Now().Add(30 * time.Second))
	c.conn.Write([]byte(GeneratePrompt(c.config.Auth.LoginPrompt)))
	login := c.authRead()

	c.conn.SetDeadline(time.Now().Add(30 * time.Second))
	c.conn.Write([]byte(GeneratePrompt(c.config.Auth.PasswordPrompt)))
	password := c.authRead()

	if c.config.Captcha.Enabled {
		captcha := c.createCaptcha()
		c.conn.Write([]byte(c.createCaptchaPrompt(captcha)))

		captchainp := c.authRead()

		if captchainp != captcha {
			c.conn.Write([]byte(GeneratePrompt(c.config.Auth.CaptchaError)))
			time.Sleep(2 * time.Second)
			c.conn.Close()
			return false
		}

	}

	if !database.CheckUser(utils.Sha3(login), utils.Sha3(password)) {
		c.conn.Write([]byte(GeneratePrompt(c.config.Auth.AuthError)))
		c.conn.Close()
		return false

	}
	c.login = login

	if c.config.Logs.NewClientConnected {
		log := c.config.Logs.NewClientConnectedFormat

		log = strings.ReplaceAll(log, "{date}", time.Now().Format("15:04:05"))
		log = strings.ReplaceAll(log, "{login}", login)
		log = strings.ReplaceAll(log, "{ip}", strings.Split(c.conn.RemoteAddr().String(), ":")[0])
		log = strings.ReplaceAll(log, "{port}", strings.Split(c.conn.RemoteAddr().String(), ":")[1])

		fmt.Println(GeneratePrompt(log))
	}

	return true
}

func (c *Connection) authRead() string {

	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)
	_data, err := tp.ReadLine()

	if err != nil {
		return ""
	}

	data := bytes.TrimPrefix([]byte(_data), []byte{255, 251, 31, 255, 251, 32, 255, 251, 24, 255, 251, 39, 255, 253, 1, 255, 251, 3, 255, 253, 3})

	return string(data)

}

func (c *Connection) createCaptcha() string {
	s1 := random.NewSource(time.Now().UnixNano())
	r1 := random.New(s1)

	var letters = []rune(c.config.Captcha.Letters)
	s := make([]rune, c.config.Captcha.CaptchaLen)

	for i := range s {
		s[i] = letters[r1.Intn(len(letters))]
	}
	return string(s)

}
