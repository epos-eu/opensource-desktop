package main

import (
	"bufio"
	"os/exec"
	"strings"
)

// Call kubeconfig to get the contexts
func (a *App) GetKubernetesContexts() ([]string, error) {
	var contexts []string

	// Run the command and get the output
	output, err := RunCommand(exec.Command("kubectl", "config", "get-contexts", "-o=name"))
	if err != nil {
		return contexts, err
	}

	// Read the output string and parse the contexts
	lines := bufio.NewScanner(strings.NewReader(output))
	for lines.Scan() {
		contexts = append(contexts, lines.Text())
	}
	if err := lines.Err(); err != nil {
		return contexts, err
	}

	return contexts, nil
}
