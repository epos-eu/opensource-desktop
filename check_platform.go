package main

import (
	"database/sql"
	"os"
	"os/exec"
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
	// command := exec.Command("docker", "compose", "--version")
	// err = command.Run()
	_, err = RunCommand(exec.Command("docker", "compose", "--version"))
	if err != nil {
		// command = exec.Command("docker-compose", "--version")
		// err = command.Run()
		_, err = RunCommand(exec.Command("docker-compose", "--version"))
		if err != nil {
			return false
		}
	}
	return true
}

func (a *App) IsDockerRunning() bool {
	// Run the command to see if docker is running
	// command := exec.Command("docker", "info")
	// _, err := command.Output()
	_, err := RunCommand(exec.Command("docker", "info"))
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
	// command := exec.Command("kubectl", "version", "--client")
	// _, err = command.Output()
	_, err = RunCommand(exec.Command("kubectl", "version", "--client"))
	return err == nil
}
