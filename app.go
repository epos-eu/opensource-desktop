package main

import (
	"bufio"
	"context"
	"fmt"
	"regexp"

	"database/sql"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	dockerMethods "github.com/epos-eu/opensource-docker/cmd/methods"
	kubernetesMethods "github.com/epos-eu/opensource-kubernetes/cmd/methods"
	_ "github.com/mattn/go-sqlite3"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

type Environment struct {
	Platform         string           `json:"platform"`
	EnvironmentSetup EnvironmentSetup `json:"environmentSetup"`
	Variables        []Section        `json:"variables"`
	AccessPoints     EposAccessPoints `json:"accessPoints"`
}

type EnvironmentSetup struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Context string `json:"context"` // Only used for kubernetes
}

type Section struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

var databasePath string

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize the database
	err := databaseInit()
	if err != nil {
		// TODO: do this in the frontend
		wailsRuntime.MessageDialog(ctx, wailsRuntime.MessageDialogOptions{
			Type:    wailsRuntime.ErrorDialog,
			Title:   "Error initializing the database",
			Message: fmt.Sprintf("Error initializing the database: %v", err),
		})
		// Exit the app
		wailsRuntime.Quit(ctx)
	}
}

// check if two environments are equal
func (e Environment) Equals(other Environment) bool {
	return e.EnvironmentSetup.Name == other.EnvironmentSetup.Name && e.EnvironmentSetup.Version == other.EnvironmentSetup.Version
}

func (a *App) GetVersion() string {
	return VERSION
}

// Checks if an environment is installed:
//   - If the platform is docker, return true if there is an environment with the same name, version and platform
//
// - If the platform is kubernetes, return true if there is an environment with the same name, version, platform and context
func (a *App) IsEnvironmentInstalled(oName, oVersion, oPlatform, oContext string) bool {
	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return false
	}
	defer db.Close()

	// Query the database for the environment
	rows, err := db.Query("SELECT name, version, platform, context FROM environments WHERE name = ? AND version = ? AND platform = ?", oName, oVersion, oPlatform)
	if err != nil {
		return false
	}
	defer rows.Close()

	// Define the variables to hold the values from the database
	var name, version, platform, context string

	// Load the values from the database
	for rows.Next() {
		err = rows.Scan(&name, &version, &platform, &context)
		if err != nil {
			return false
		}

		// If the platform is kubernetes, check if the context is the same
		if platform == "docker" && platform == oPlatform && name == oName && version == oVersion {
			return true
		}

		// If the platform is kubernetes, check if the context is the same
		if platform == "kubernetes" && platform == oPlatform && name == oName && version == oVersion && context == oContext {
			return true
		}
	}
	return false
}

func getEnvironmentVariablesTempFilePath(envName, envVersion, platform string) (string, error) {
	// Get the variables
	environment, err := getInstalledEnvironment(envName, envVersion, platform)
	if err != nil {
		return "", err
	}
	variables := environment.Variables

	// Generate the env file
	envTempFilePath, err := generateTempFile(os.TempDir(), "env", variablesToBinary(variables))
	return envTempFilePath, err
}

func (a *App) GetInstalledEnvironments() ([]Environment, error) {
	// TODO use the getInstalledEnvironment function
	// Return all the installed environments from the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database for all the environments
	rows, err := db.Query("SELECT name, version, platform, variables, context, apiGateway, dataPortal FROM environments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the environments
	var environments []Environment

	// Iterate over the rows and add them to the slice
	for rows.Next() {
		var name, version, platform, variables, context, apiGateway, dataPortal string
		err = rows.Scan(&name, &version, &platform, &variables, &context, &apiGateway, &dataPortal)
		if err != nil {
			return nil, err
		}

		// Convert the variables to a slice of Section
		var sections []Section
		err = json.Unmarshal([]byte(variables), &sections)
		if err != nil {
			return nil, err
		}

		// Add the environment to the slice
		environments = append(environments, Environment{
			Platform:         platform,
			EnvironmentSetup: EnvironmentSetup{Name: name, Version: version, Context: context},
			Variables:        sections,
			AccessPoints:     EposAccessPoints{ApiGateway: apiGateway, DataPortal: dataPortal},
		})
	}

	// For each environment, check if it is still installed
	for i, environment := range environments {

		if environment.Platform == "docker" {
			// Get the installed docker environments from the docker ps command
			output, err := RunCommand(exec.Command("docker", "ps", "-a", "--format", "{{.Names}}"))
			if err != nil {
				return nil, err
			}

			// Get the tagname of the environment
			envName := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(environment.EnvironmentSetup.Name+environment.EnvironmentSetup.Version, "-")

			// Check if the tagname is in the output of the command
			if !strings.Contains(output, envName) {
				// If it is not, that means that the environment is not installed in the system anymore (it was removed manually)
				// Remove the environment from the slice
				environments = append(environments[:i], environments[i+1:]...)
				// Remove the environment from the database
				err = deleteEnvironmentFromDatabase(environment.EnvironmentSetup.Name, environment.EnvironmentSetup.Version, environment.Platform, environment.EnvironmentSetup.Context)
				if err != nil {
					return nil, err
				}
			}
		}

		if environment.Platform == "kubernetes" {
			// kubectl config use-context <context>
			_, err = RunCommand(exec.Command("kubectl", "config", "use-context", environment.EnvironmentSetup.Context))
			if err != nil {
				// If the context is not valid, remove the environment from the slice
				environments = append(environments[:i], environments[i+1:]...)
				// Remove the environment from the database
				err = deleteEnvironmentFromDatabase(environment.EnvironmentSetup.Name, environment.EnvironmentSetup.Version, environment.Platform, environment.EnvironmentSetup.Context)
				if err != nil {
					return nil, err
				}
			}

			// kubectl get namespaces
			output, err := RunCommand(exec.Command("kubectl", "get", "namespaces", "--no-headers", "-o", "custom-columns=NAME:.metadata.name"))
			if err != nil {
				return nil, err
			}

			// Check if the namespace is in the output of the command
			if !strings.Contains(output, environment.EnvironmentSetup.Name) {
				// If it is not, that means that the environment is not installed in the system anymore (it was removed manually)
				// Remove the environment from the slice
				environments = append(environments[:i], environments[i+1:]...)
				// Remove the environment from the database
				err = deleteEnvironmentFromDatabase(environment.EnvironmentSetup.Name, environment.EnvironmentSetup.Version, environment.Platform, environment.EnvironmentSetup.Context)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Sort the environments by name and version
	sort.Slice(environments, func(i, j int) bool {
		if environments[i].EnvironmentSetup.Name == environments[j].EnvironmentSetup.Name {
			return environments[i].EnvironmentSetup.Version > environments[j].EnvironmentSetup.Version
		}
		return environments[i].EnvironmentSetup.Name < environments[j].EnvironmentSetup.Name
	})

	return environments, nil
}

func getInstalledEnvironment(name, version, platform string) (Environment, error) {
	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return Environment{}, err
	}
	defer db.Close()

	// Query the database for the environment
	rows, err := db.Query("SELECT variables, context, apiGateway, dataPortal FROM environments WHERE name = ? AND version = ?", name, version, platform)
	if err != nil {
		return Environment{}, err
	}
	defer rows.Close()

	// If the query didn't match any rows, return an error
	if !rows.Next() {
		return Environment{}, fmt.Errorf("environment not found: %s %s", name, version)
	}

	// Get the variables from the database
	var variables, context, apiGateway, dataPortal string
	err = rows.Scan(&variables, &context, &apiGateway, &dataPortal)
	if err != nil {
		return Environment{}, err
	}

	// Convert the variables to a slice of Section
	var sections []Section
	err = json.Unmarshal([]byte(variables), &sections)
	if err != nil {
		return Environment{}, err
	}

	return Environment{
		Platform:         platform,
		EnvironmentSetup: EnvironmentSetup{Name: name, Version: version, Context: context},
		Variables:        sections,
		AccessPoints:     EposAccessPoints{ApiGateway: apiGateway, DataPortal: dataPortal},
	}, nil
}

// Check if there is an internet connection if there isn't, show a message dialog and exit
// TODO: do this in the frontend
func (a *App) IsInternetConnected() bool {
	_, err := net.DialTimeout("tcp", "www.google.com:80", time.Second*2)
	if err != nil {
		// Show a message dialog
		wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
			Type:    wailsRuntime.ErrorDialog,
			Title:   "Error",
			Message: "No internet connection. An internet connection is required to install a new environment, please try again later.",
		})
		// Exit the app
		wailsRuntime.Quit(a.ctx)
	}
	return true
}

// Gets the ip address of the machine
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func (a *App) GetIp() (string, error) {
	// Dial a UDP connection to a public IP address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	// Ensure the connection is closed
	defer conn.Close()

	// Get the local network address from the connection.
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// Return the IP address as a string.
	return localAddr.IP.String(), nil
}

// Read the env.env file for the given platform and return the sections with their variables
func (a *App) ReadEnvVariables(platform string) ([]Section, error) {
	// Initialize the variables slice and the error
	var variables []Section
	err := fmt.Errorf("platform not supported: %s", platform)

	// Read the env file for the given platform
	if platform == "docker" {
		// Read the env file used by the docker cmd
		variables, err = readEnvFile(dockerMethods.GetConfigurationsEmbed())
	} else if platform == "kubernetes" {
		// Read the env file used by the kubernetes cmd
		variables, err = readEnvFile(kubernetesMethods.GetConfigurationsEmbed())
	}
	return variables, err
}

// Parse the env.env file and return the sections with their variables
func readEnvFile(file []byte) ([]Section, error) {
	// Parse the env file
	var sections []Section
	var currentSection *Section
	lines := strings.Split(string(file), "\n")

	delimiter := "# ************************************************************************************************************"

	// Iterate over the lines and parse the sections and variables
	for i, line := range lines {
		// Check if the line is a section header:
		// 		# ************************************************************************************************************
		// 		#                                       SOME TEXT
		// 		# ************************************************************************************************************
		if strings.HasPrefix(line, "#") {
			// Check if the line before and after follows the pattern
			// Before
			if i > 0 && lines[i-1] == delimiter {
				// After
				if i < len(lines)-1 && lines[i+1] == delimiter {
					// Create a new section
					currentSection = &Section{
						// Remove the leading and trailing spaces and the # character
						Name:      strings.TrimSpace(line)[1:],
						Variables: make(map[string]string),
					}
					sections = append(sections, *currentSection)
				}
			}
			// Else if the line is a variable (not a comment)
		} else if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			// Parse the variable and its default value
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				// Skip lines that don't look like key=value
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove the quotes from the value
			value = strings.Trim(value, "\"")

			// Add the variable to the current section
			if currentSection != nil {
				currentSection.Variables[key] = value
			}
		}
	}

	return sections, nil
}

// Initialize the database
func databaseInit() error {
	// TODO: see where to put the database on each platform
	// dbPath := "./environments.db"
	dbPath, err := getDatabasePath()
	if err != nil {
		return err
	}
	databasePath = dbPath + "environments.db"

	// Create the database file
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the environments table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS environments (
        name TEXT,
        version TEXT,
        platform TEXT,
        context TEXT,
        dataPortal TEXT,
        apiGateway TEXT,
        variables TEXT,
        PRIMARY KEY (name, version, platform)
    );

	CREATE TABLE IF NOT EXISTS platform_paths (
		platform TEXT,
		path TEXT,
		PRIMARY KEY (platform)
	);
    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

// Get the path to the folder where to save the database
func getDatabasePath() (string, error) {
	folder := "EPOS_opensource_desktop"
	var basePath string

	// If on windows, use the appdata folder
	if runtime.GOOS == "windows" {
		basePath = filepath.Join(os.Getenv("APPDATA"), folder)
	}

	// If on linux, use the home folder
	if runtime.GOOS == "linux" {
		basePath = filepath.Join(os.Getenv("HOME"), folder)
	}

	// If on mac, use the home folder
	if runtime.GOOS == "darwin" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		// Go in Application Support
		basePath = filepath.Join(home, "Library", "Application Support", folder)
	}

	if basePath == "" {
		return "", fmt.Errorf("unsupported operating system")
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0755)
		if err != nil {
			return "", err
		}
	}

	// Return the path to the database file within the directory
	return basePath + string(filepath.Separator), nil
}

// Open a file dialog to select a path for the location of the platform
func (a *App) SpecifyPlatformPath(platform string) (string, error) {
	title := "Select the location of the " + platform + " installation"
	path, err := a.OpenFolderDialog(title)
	if err != nil {
		return "", err
	}

	// Save the path to the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Update the database with the path
	_, err = db.Exec("INSERT OR REPLACE INTO platform_paths(platform, path) VALUES(?, ?)", platform, path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// Open the file dialog to select a folder
func (a *App) OpenFolderDialog(title string) (string, error) {
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:                title,
		CanCreateDirectories: false,
		ShowHiddenFiles:      true,
	})
	if err != nil {
		return "", err
	}

	// Add a trailing slash to the path
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}

	return path, nil
}

// Check if the port is available
func (a *App) IsPortAvailable(port string) (bool, error) {
	// Validate the string
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		return false, err
	}

	// Check if the port is used by the environments
	usedPorts, err := getUsedPorts()
	if err != nil {
		return false, err
	}
	for _, usedPort := range usedPorts {
		if usedPort == port {
			return false, nil
		}
	}

	// Try to listen on the port
	ln, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return true, nil
	}
	defer ln.Close()

	return false, nil
}

// Get the available port
func (a *App) GetAvailablePort() (string, error) {
	// TODO: make this more efficient
	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		// Try to listen on the port, 0 means a port number is automatically chosen
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return "", err
		}
		defer ln.Close()

		// Get the port from the listener
		addr := ln.Addr().String() // "ip:port"
		parts := strings.Split(addr, ":")
		port := parts[len(parts)-1]

		// Check if the port is disallowed
		if !isPortDisallowed(port) {
			return port, nil
		}
	}

	return "", fmt.Errorf("could not find an available port")
}

func isPortDisallowed(port string) bool {
	// Get the used ports
	usedPorts, err := getUsedPorts()
	if err != nil {
		return true
	}

	// Check if the port is disallowed
	for _, usedPort := range usedPorts {
		if usedPort == port {
			return true
		}
	}

	return false
}

// Get all the ports already used by the environments
func getUsedPorts() ([]string, error) {
	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database for the variables
	rows, err := db.Query("SELECT variables FROM environments")
	if err != nil {
		return nil, err
	}

	// Parse the variables and get the ports
	var usedPorts []string
	for rows.Next() {
		var variables string
		err = rows.Scan(&variables)
		if err != nil {
			return nil, err
		}

		// Convert the variables to a slice of Section
		var sections []Section
		err = json.Unmarshal([]byte(variables), &sections)
		if err != nil {
			return nil, err
		}

		// Iterate over the sections and get the ports
		for _, section := range sections {
			for key, value := range section.Variables {
				if strings.Contains(key, "_PORT") {
					usedPorts = append(usedPorts, value)
				}
			}
		}
	}

	return usedPorts, nil
}

// Call the docker cmd to populate an environment
func (a *App) PopulateEnvironment(envName, envTag, path, platform string) error {
	// Create a temporary file with the environment variables
	envFilePath, err := getEnvironmentVariablesTempFilePath(envName, envTag, platform)
	if err != nil {
		return err
	}

	// Log the parameters just for debugging
	fmt.Println("path: ", path)
	fmt.Println("envName: ", envName)
	fmt.Println("envTag: ", envTag)
	fmt.Println("envFilePath: ", envFilePath)
	fmt.Println("platform: ", platform)

	// Intercept the output of the command
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to wait for the command to finish
	done := make(chan error)

	// Run the populate in a goroutine
	go func() {
		if platform == "docker" {
			// Run the command and get the error
			err = dockerMethods.PopulateEnvironment(
				envFilePath, // environment variables file path
				path,        // path to the environment
				envName,     // environment name
				envTag,      // environment tag
			)
		} else if platform == "kubernetes" {
			// I have to get the context from the database because the frontend doesn't have it when calling this function (should probably be fixed in the frontend)
			context, err := getKubernetesEnvironmentContext(envName, envTag)
			if err != nil {
				// Restore the real stdout
				w.Close()
				os.Stdout = old
				done <- err
				return
			}

			// Run the command and get the error
			err = kubernetesMethods.PopulateEnvironment(
				context,     // kubernetes context
				envFilePath, // environment variables file path
				path,        // path to the files to populate
				envName,     // environment name (namespace)
				envTag,      // environment tag
			)
		} else {
			// This should never happen
			err = fmt.Errorf("unknown platform: %s", platform)
		}

		// Restore the real stdout
		w.Close()
		os.Stdout = old

		// Send the error to the channel
		done <- err
	}()

	// Create a scanner to read the output line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Emit the events to the frontend for each line
		wailsRuntime.EventsEmit(a.ctx, "TERMINAL_OUTPUT", scanner.Text())
	}

	// Wait for the command to finish
	err = <-done

	//Remove the temporary file even if there was an error
	os.Remove(envFilePath)

	return err
}

// Get the context for a kubernetes environment from the db
func getKubernetesEnvironmentContext(envName, envVersion string) (string, error) {
	// Open the database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Query the database for the context
	rows, err := db.Query("SELECT context FROM environments WHERE name = ? AND version = ? AND platform = ?", envName, envVersion, "kubernetes")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Load the context from the database
	context := ""
	for rows.Next() {
		err = rows.Scan(&context)
		if err != nil {
			return "", err
		}
	}

	// If the context is empty, return an error
	if context == "" {
		return "", fmt.Errorf("context not found: %s %s", envName, envVersion)
	}

	return context, nil
}
