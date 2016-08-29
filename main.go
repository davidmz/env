package main // import "gopkg.in/davidmz/env.v1"

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	args := os.Args
	args = args[1:]

	currentEnv := allEnvVars()

	if len(args) > 0 && args[0] == "-i" {
		currentEnv = nil
		args = args[1:]
	}

	for len(args) > 0 {
		p := strings.SplitN(args[0], "=", 2)
		if len(p) != 2 || p[0] == "" {
			break
		}
		currentEnv.Set(p[0], p[1])
		args = args[1:]
	}

	if len(args) == 0 {
		for _, e := range currentEnv {
			fmt.Println(e)
		}
		return
	}

	exe := args[0]
	args = args[1:]
	exePath, err := exec.LookPath(exe)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(exePath, args...)
	cmd.Env = currentEnv.ToStrings()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			os.Exit(ws.ExitStatus())
		} else {
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func allEnvVars() (res envVars) {
	for _, e := range os.Environ() {
		p := strings.SplitN(e, "=", 2)
		// skip special windows variables started from '='
		// http://stackoverflow.com/q/30102750
		if p[0] != "" {
			res = append(res, envVar{p[0], p[1]})
		}
	}
	return
}

type envVar struct {
	name  string
	value string
}

func (e envVar) String() string { return e.name + "=" + e.value }

type envVars []envVar

func (vars *envVars) Set(name, value string) {
	for i, e := range *vars {
		if e.name == name {
			(*vars)[i].value = value
			return
		}
	}
	(*vars) = append((*vars), envVar{name, value})
}

func (vars envVars) ToStrings() (res []string) {
	for _, e := range vars {
		res = append(res, e.String())
	}
	return
}
