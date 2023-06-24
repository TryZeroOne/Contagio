package cnc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	USER_CONNECTED int = 1
	ATTACK_STARTED int = 2
)

type Logging struct {
	action   int
	c        *Connection
	target   string
	port     string
	duration string
	method   string
}

func (l *Logging) SendLog() {
	defer Catch()

	config := l.c.config.Logs

	switch l.action {
	// crazy ahaha
	case USER_CONNECTED:
		l.c.SendTelegramLog(config.NewClientConnectedTelegram, l)
		l.c.PrintLog(config.NewClientConnectedTerminal, l)
		l.c.SaveLog(config.NewClientConnectedFile, l, config.NewClientConnectedFileName)
	case ATTACK_STARTED:
		l.c.SendTelegramLog(config.NewAttackStartedTelegram, l)
		l.c.SaveLog(config.NewAttackStartedFile, l, config.NewAttackStartedFileName)
		l.c.PrintLog(config.NewAttackStartedTerminal, l)
	}

}

func (l *Logging) formatLog(log string) string {
	defer Catch()

	log = strings.ReplaceAll(log, "{date}", time.Now().Format("15:04:05"))
	log = strings.ReplaceAll(log, "{login}", l.c.login)
	log = strings.ReplaceAll(log, "{ip}", strings.Split(l.c.conn.RemoteAddr().String(), ":")[0])
	log = strings.ReplaceAll(log, "{port}", strings.Split(l.c.conn.RemoteAddr().String(), ":")[1])

	if l.action == int(ATTACK_STARTED) {
		log = strings.ReplaceAll(log, "{method}", l.method)
		log = strings.ReplaceAll(log, "{target}", l.target)
		log = strings.ReplaceAll(log, "{duration}", l.duration)
		log = strings.ReplaceAll(log, "{target_port}", l.port)
	}

	return log
}

func (c *Connection) SaveLog(log string, l *Logging, filename string) {
	defer Catch()

	config := c.config.Logs
	if config.SaveLogsInFile {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Can't open: " + err.Error())
			return
		}
		defer file.Close()

		_, err = file.Write([]byte(l.formatLog(log) + "\n"))
		if err != nil {
			fmt.Println("Can't write: " + err.Error())
		}

	}
}

func (c *Connection) PrintLog(log string, l *Logging) {
	defer Catch()

	if c.config.Logs.PrintLogsInTerminal {
		fmt.Println(GeneratePrompt(l.formatLog(log)))
	}
}

func (c *Connection) SendTelegramLog(log string, l *Logging) {
	defer Catch()

	config := c.config.Logs
	if config.SendLogsInTelegram {
		log = l.formatLog(log) + "*"
		text := strings.ReplaceAll(log, "\n", "%0D%0A")
		text = strings.ReplaceAll(text, ".", "\\.")
		telegram, _ := http.Get("https://api.telegram.org/bot" + config.TelegramBotToken + "/sendMessage?chat_id=" + config.TelegramChatId + "&parse_mode=MarkdownV2&text=" + text + "*")

		if telegram.StatusCode != 200 {
			body, err := io.ReadAll(telegram.Body)
			if err != nil {
				return
			}
			fmt.Println("Unknown error: " + string(body))
		}
	}
}

func (c *Connection) NewLog(action int, target, method, duration, port string) *Logging {
	defer Catch()

	return &Logging{
		action:   action,
		target:   target,
		port:     port,
		duration: duration,
		method:   method,
		c:        c,
	}

}
