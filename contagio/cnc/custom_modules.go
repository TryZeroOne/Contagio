package cnc

import (
	"os"
	"os/exec"
	"strings"
)

func (c *Connection) FormatModule(module string) (result string) {

	var tempmodule = module

	for name := range c.config.Modules {
		tempmodule = strings.ReplaceAll(tempmodule, "{"+name+"}", run(c.config.Modules[name].Exec, c.config.Modules[name].ExecEnv, c.config.Modules[name].ExecDir))
	}

	return tempmodule
}

func run(command, exenv, dir string) (output string) {

	parts := strings.Split(command, " ")
	var execdir string
	var err error
	if dir == "" {
		execdir, err = os.Getwd()
		if err != nil {
			return err.Error()
		}
	} else {
		execdir = dir
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = execdir
	var out strings.Builder

	env := os.Environ()

	if exenv != "" {
		env = append(env, exenv)
		cmd.Env = env
	}

	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return err.Error()
	}
	return out.String()

}
