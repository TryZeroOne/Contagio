package cnc

import (
	"contagio/contagio/bot_server"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func (c *Connection) CncMainMenu() {

	c.Cls()

	c.conn.SetDeadline(time.Now().Add(time.Duration(1 * time.Hour)))

	for _, i := range c.config.Cnc.Banner {
		_, err := c.conn.Write([]byte(GeneratePrompt(i)))
		if err != nil {
			c.conn.Close()
			return
		}
	}

	for {
		err := c.CommandHandler()
		if err != nil {
			return
		}
	}

}

func (c *Connection) GenerateCmdPrompt(str string) string {

	prompt := GeneratePrompt(str)

	prompt = strings.ReplaceAll(prompt, "{login}", c.login)

	return prompt

}

func (c *Connection) Title() {
	for {
		count := bot_server.BotCount

		title := c.GenerateCmdPrompt(c.config.Cnc.Title)

		title = strings.ReplaceAll(title, "{bots}", strconv.Itoa(count))
		title = strings.ReplaceAll(title, "{memory}", strconv.Itoa(GetMemoryUsage()))
		title = strings.ReplaceAll(title, "{cpu}", strconv.Itoa(int(cpuUsage())))
		title = c.FormatModule(title)

		if c.config.Animation.Enabled {
			var animtext = c.config.Animation.Letters
			for i := 0; i < len(animtext); i++ {
				c.conn.Write([]byte("\033]0;" + strings.ReplaceAll(title, "{animation}", animtext[:i+1]) + "\007"))

				time.Sleep(time.Duration(c.config.Animation.Delay) * time.Millisecond)
			}
			for i := len(animtext) - 1; i >= -1; i-- {
				c.conn.Write([]byte("\033]0;" + strings.ReplaceAll(title, "{animation}", animtext[:i+1]) + "\007"))
				time.Sleep(time.Duration(c.config.Animation.Delay) * time.Millisecond)
			}

		} else {
			c.conn.Write([]byte("\033]0;" + title + "\007"))

		}

	}

}

func cpuUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if int(percent[0]) == 0 {
		percent[0] = 1
	}
	return percent[0]
}

func GetMemoryUsage() int {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return int(math.Ceil(memory.UsedPercent))
}
