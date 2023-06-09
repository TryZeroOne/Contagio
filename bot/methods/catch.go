package methods

import "fmt"

func Catch() {
	if er := recover(); er != nil {
		fmt.Print(er)
		return
	}
}
