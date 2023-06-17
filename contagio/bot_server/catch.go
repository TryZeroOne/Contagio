package bot_server

import "fmt"

func catch() {
	if err := recover(); err != nil {
		fmt.Println("[contagio] Fatal error: " + err.(string))
	}
}
