package main

import (
	"bufio"

	dockerMethods "github.com/epos-eu/opensource-docker/cmd/methods"
	kubernetesMethods "github.com/epos-eu/opensource-kubernetes/cmd/methods"

	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// dockerCmd "github.com/epos-eu/opensource-docker/cmd"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	_ "github.com/mattn/go-sqlite3"
)

type EposAccessPoints struct {
	ApiGateway string `json:"apiGateway"`
	DataPortal string `json:"dataPortal"`
}

func (a *App) InstallEnvironment(platform string, environmentSetup EnvironmentSetup, variables []Section, skipImagesAutoupdate bool, isEdit bool) error {

	// print for debugging
	fmt.Println("Platform: ", platform)
	fmt.Println("EnvironmentSetup: ", environmentSetup)
	fmt.Println("Variables: ", variables)
	fmt.Println("SkipImagesAutoupdate: ", skipImagesAutoupdate)
	fmt.Println("IsEdit: ", isEdit)

	var err error
	var accessPoints EposAccessPoints

	// TODO run the install script
	if platform == "docker" {
		err = a.installDockerEnvironment(environmentSetup, variables, skipImagesAutoupdate, isEdit)

		// Build the access points strings
		accessPoints = EposAccessPoints{
			DataPortal: "http://" + os.Getenv("API_HOST_ENV") + ":" + os.Getenv("DATA_PORTAL_PORT"),
			ApiGateway: "http://" + os.Getenv("API_HOST_ENV") + ":" + os.Getenv("API_PORT") + os.Getenv("DEPLOY_PATH") + os.Getenv("API_PATH") + "/ui/",
		}
	} else if platform == "kubernetes" {
		err = a.installKubernetesEnvironment(environmentSetup, variables, skipImagesAutoupdate, isEdit)

		// Build the access points strings
		accessPoints = EposAccessPoints{
			DataPortal: os.Getenv("PORTAL_URL_READY"),
			ApiGateway: os.Getenv("API_URL_READY"),
		}
	} else {
		return fmt.Errorf("unknown platform: %s", platform)
	}
	if err != nil {
		return err
	}

	// Save the environment to the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Convert the variables to a JSON string
	variablesJson, err := json.Marshal(variables)
	if err != nil {
		return err
	}

	// TODO: maybe get the ports from the environment variables istead of using the values from the variables variable (the deploy might change them if they are already in use)

	// Upsert the environment into the database
	_, err = db.Exec("INSERT OR REPLACE INTO environments(name, version, platform, dataPortal, apiGateway, variables, context) VALUES(?, ?, ?, ?, ?, ?, ?)",
		environmentSetup.Name,
		environmentSetup.Version,
		platform,
		accessPoints.DataPortal,
		accessPoints.ApiGateway,
		string(variablesJson),
		environmentSetup.Context,
	)
	if err != nil {
		return err
	}

	// Return nil if there was no error
	return nil
}

func (a *App) installDockerEnvironment(environmentSetup EnvironmentSetup, variables []Section, skipImagesAutoupdate bool, isEdit bool) error {
	// Generate a temporary file with the environment variables
	envTempFilePath, err := generateTempFile(os.TempDir(), "configurations", variablesToBinary(variables))
	if err != nil {
		return err
	}

	// Add the arguments to the docker command
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to wait for the command to finish
	done := make(chan error)

	go func() {
		// Run the docker command
		err := dockerMethods.CreateEnvironment(
			envTempFilePath,                         // the file with the environment variables
			"",                                      // the docker-compose file
			"",                                      // external ip
			environmentSetup.Name,                   // the name of the environment
			environmentSetup.Version,                // the version of the environment
			fmt.Sprintf("%t", isEdit),               // if the environment is being edited/updated
			fmt.Sprintf("%t", skipImagesAutoupdate), // if the images should be updated
		)

		// back to normal state
		w.Close()
		os.Stdout = old // restoring the real stdout
		done <- err
	}()

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Emit the events to the frontend for each line
		wailsRuntime.EventsEmit(a.ctx, "TERMINAL_OUTPUT", scanner.Text())
	}

	err = <-done // wait for the command to finish

	//Remove the temporary file
	os.Remove(envTempFilePath)

	return err
}

func (a *App) installKubernetesEnvironment(environmentSetup EnvironmentSetup, variables []Section, skipImagesAutoupdate bool, isEdit bool) error {
	// Generate a temporary file with the environment variables
	envTempFilePath, err := generateTempFile(os.TempDir(), "configurations", variablesToBinary(variables))
	if err != nil {
		return err
	}

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to wait for the command to finish
	done := make(chan error)

	go func() {
		// Run the kubernetes command
		err := kubernetesMethods.CreateEnvironment(
			envTempFilePath,                         // the file with the environment variables
			environmentSetup.Context,                // the context
			environmentSetup.Name,                   // the namespace
			environmentSetup.Version,                // the version of the environment
			fmt.Sprintf("%t", skipImagesAutoupdate), // if the images should be updated
			fmt.Sprintf("%t", isEdit),               // if the environment is being edited/updated
		)

		w.Close()
		os.Stdout = old // restoring the real stdout
		done <- err
	}()

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Emit the events to the frontend for each line
		wailsRuntime.EventsEmit(a.ctx, "TERMINAL_OUTPUT", scanner.Text())
	}

	err = <-done // wait for the command to finish

	//Remove the temporary file
	os.Remove(envTempFilePath)

	return err
}

// Convert the variables to a binary to be saved in a file
func variablesToBinary(variables []Section) []byte {
	var result []string
	for _, section := range variables {
		variables := section.Variables // a map of variable names to values
		for name, value := range variables {
			result = append(result, fmt.Sprintf("%s=%s", name, value))
		}
	}

	return []byte(strings.Join(result, "\n"))
}

// Generate a temporary file with the given data and return the file path
func generateTempFile(dname string, filetype string, text []byte) (string, error) {
	// If dname already exists, remove it
	if _, err := os.Stat(dname); err == nil {
		os.RemoveAll(dname)
	}

	tmpFile, err := os.CreateTemp(dname, filetype)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	name := tmpFile.Name()
	if _, err = tmpFile.Write(text); err != nil {
		return "", err
	}

	return name, nil
}
