package main

import (
	"database/sql"
	"fmt"

	dockerMethods "github.com/epos-eu/opensource-docker/cmd/methods"
	kubernetesMethods "github.com/epos-eu/opensource-kubernetes/cmd/methods"
)

// Deletes an installed environment from the database given its name and version
func (a *App) DeleteInstalledEnvironment(platform, name, version, context string) error {

	fmt.Println("Platform: ", platform)
	fmt.Println("Name: ", name)
	fmt.Println("Version: ", version)
	fmt.Println("Context: ", context)

	var err error

	if platform == "docker" {
		err = deleteDockerEnvironment(name, version)
	} else if platform == "kubernetes" {
		err = deleteKubernetesEnvironment(name, context)
	} else {
		return fmt.Errorf("unknown platform: %s", platform)
	}
	if err != nil {
		return err
	}

	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the environment from the database
	result, err := db.Exec("DELETE FROM environments WHERE name = ? AND version = ? AND platform = ?", name, version, platform)
	// If there was an error or the query didn't match any rows, return the error
	if err != nil || result == nil {
		return err
	}

	return nil
}

func deleteDockerEnvironment(name, version string) error {
	// Get the environment variables as a temp file
	envFilePath, err := getEnvironmentVariablesTempFilePath(name, version, "docker")

	// Call the delete cmd
	err = dockerMethods.DeleteEnvironment(
		envFilePath, // environment variables file path
		"",          // docker compose file path
		name,        // environment name
		version,     // environment version
	)
	return err
}

func deleteKubernetesEnvironment(name, context string) error {
	// Call the delete cmd
	return kubernetesMethods.DeleteEnvironment(
		context, // kubernetes context
		name,    // namespace
	)
}
