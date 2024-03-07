//go:build windows
// +build windows

package main

import (
	"database/sql"
	"os"
	"os/exec"
	"syscall"
)

// See if Docker is installed
func (a *App) IsDockerInstalled() bool {
	// Add +"usr/local/bin:" to the PATH
	os.Setenv("PATH", "/usr/local/bin:"+os.Getenv("PATH"))

	// Add to the PATH the location of the docker and docker-compose executables from the database
	db, err := sql.Open("sqlite3", databasePath)
	if err == nil {
		defer db.Close()
		// Query the database for the path to the docker and docker-compose executables
		rows, err := db.Query("SELECT path FROM platform_paths WHERE platform = ?", "docker")
		if err == nil {
			defer rows.Close()
			// If the query matched a row, add the path to the PATH
			if rows.Next() {
				var path string
				if err = rows.Scan(&path); err == nil {
					os.Setenv("PATH", path+":"+os.Getenv("PATH"))
				}
			}
		}
	}

	// Check if docker is installed
	// Run "docker compose --version", if it fails, run "docker-compose --version"
	command := exec.Command("docker", "compose", "--version")
	// This is only needed on Windows to hide the console window.
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = command.Run()
	if err != nil {
		command = exec.Command("docker-compose", "--version")
		// This is only needed on Windows to hide the console window.
		command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err = command.Run()
		if err != nil {
			return false
		}
	}
	return true
}

func (a *App) IsDockerRunning() bool {
	// Run the command to see if docker is running
	command := exec.Command("docker", "info")
	// This is only needed on Windows to hide the console window.
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := command.Output()
	return err == nil
}

// See if Kubernetes is installed
func (a *App) IsKubernetesInstalled() bool {
	// Add +"usr/local/bin:" to the PATH
	os.Setenv("PATH", "/usr/local/bin:"+os.Getenv("PATH"))

	// Add to the PATH the location of the kubectl executable from the database
	db, err := sql.Open("sqlite3", databasePath)
	if err == nil {
		defer db.Close()
		// Query the database for the path to the kubectl executable
		rows, err := db.Query("SELECT path FROM platform_paths WHERE platform = ?", "kubernetes")
		if err == nil {
			defer rows.Close()
			// If the query matched a row, add the path to the PATH
			if rows.Next() {
				var path string
				if err = rows.Scan(&path); err == nil {
					os.Setenv("PATH", path+":"+os.Getenv("PATH"))
				}
			}
		}
	}

	// Run the command to see if kubectl is installed
	command := exec.Command("kubectl", "version", "--client")
	// This is only needed on Windows to hide the console window.
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = command.Output()
	return err == nil
}
