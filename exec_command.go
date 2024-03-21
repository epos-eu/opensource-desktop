//go:build !windows
// +build !windows

package main

import (
	"os/exec"
)

// Run a command and return the output as a string
func RunCommand(cmd *exec.Cmd) (string, error) {
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
