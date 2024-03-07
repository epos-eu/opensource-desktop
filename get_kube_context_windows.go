//go:build windows
// +build windows

package main

import (
	"bufio"
	"os/exec"
	"syscall"
)

// Call kubeconfig to get the contexts
func (a *App) GetKubernetesContexts() ([]string, error) {
	var contexts []string

	// Run the command
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o=name")
	// This is only needed on Windows to hide the console window.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Get the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return contexts, err
	}
	if err := cmd.Start(); err != nil {
		return contexts, err
	}

	// Read the output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		contexts = append(contexts, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return contexts, err
	}

	return contexts, nil
}
