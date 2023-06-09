package cnc

import "fmt"

func Catch() {
	if err := recover(); err != nil {
		fmt.Println("[contagio] Fatal error: ", err)
		return
	}
}
