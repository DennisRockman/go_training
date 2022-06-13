package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdPrepare := exec.Command(cmd[0], cmd[1:]...) // nolint:gosec
	for key, envValue := range env {
		if err := os.Unsetenv(key); err != nil {
			return 3
		}
		if !envValue.NeedRemove {
			if err := os.Setenv(key, envValue.Value); err != nil {
				return 2
			}
		}
	}
	cmdPrepare.Env = os.Environ()
	cmdPrepare.Stdout = os.Stdout
	cmdPrepare.Stderr = os.Stderr
	if err := cmdPrepare.Run(); err != nil {
		fmt.Fprintf(cmdPrepare.Stderr, "error: %v\n", err)
		return 1
	}
	return 0
}
