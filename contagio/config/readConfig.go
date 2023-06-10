package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
)

var Release bool

type Config struct {
	Wg *sync.WaitGroup

	ImportTheme  string
	CncServer    string
	BotServer    string
	TftpServer   string
	FtpServer    string
	LoaderServer string
	RELEASE_MODE bool

	Payload struct {
		FtpPassword string
		FtpLogin    string
	}

	RootLogin string
	Logs      struct {
		NewClientConnected       bool
		NewClientConnectedFormat string
	}
	Cnc struct {
		HelpCommand string
		CmdPrompt   string

		MethodsList          string
		CustomMethods        []string
		CustomMethodsEnabled bool

		CustomHelp                []string
		CustomHelpEnabled         bool
		Banner                    []string
		InvalidCommandSyntaxError string
		UnknownCommandError       string
		BotCount                  string
		NoBotsConnectedError      string

		CommandSent string
		Title       string

		CommandExecuted      string
		CommandInvalidSyntax string
	}

	Auth struct {
		LoginPrompt         string
		PasswordPrompt      string
		AuthError           string
		CaptchaPrompt       string
		CaptchaError        string
		AllowAllIps         bool
		IpIsNotAllowedError string
	}

	Captcha struct {
		Enabled    bool
		CaptchaLen int
		Letters    string
	}

	Animation struct {
		Enabled bool
		Delay   int
		Letters string
	}

	Modules map[string]Cmodule
}

type Cmodule struct {
	Exec    string
	ExecEnv string
	ExecDir string
}

func ReadConfig(wg *sync.WaitGroup) *Config {

	var config Config

	conf, err := os.ReadFile("./config.toml")
	if err != nil {
		fmt.Printf("Config error: %e", err)
		return nil
	}
	_, err = toml.Decode(string(conf), &config)

	if err != nil {
		fmt.Printf("Config error: %e", err)
		return nil
	}

	theme, err := os.ReadFile(config.ImportTheme)

	if err != nil {
		fmt.Printf("Config error: %e", err)
		return nil
	}
	_, err = toml.Decode(string(theme), &config)

	if err != nil {
		fmt.Printf("Config error: %e", err)
		return nil
	}

	Release = config.RELEASE_MODE
	config.Wg = wg

	return &config

}
