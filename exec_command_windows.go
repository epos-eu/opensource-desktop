//go:build windows
// +build windows

package main

import (
	"os/exec"
	"syscall"
)

// Run a command with the console window hidden and return the output as a string
func RunCommand(cmd *exec.Cmd) (string, error) {
	// Only needed on Windows to hide the console window
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), err
}
