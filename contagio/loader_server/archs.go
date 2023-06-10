package loader_server

import (
	"fmt"
	"os"
	"sync"
)

var Archs = make([]string, 0)

func GetArchs(wg *sync.WaitGroup) {
	archs, err := os.ReadDir("./bin/")
	if err != nil {
		fmt.Println("[contagio] Can't read ./bin dir")
		wg.Done()
		return
	}

	for _, arch := range archs {
		Archs = append(Archs, arch.Name())
	}

}
