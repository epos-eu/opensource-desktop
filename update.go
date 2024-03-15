package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/minio/selfupdate"

	"github.com/google/go-github/v60/github"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	VERSION    = "0.0.2"
	PUBLIC_KEY = `untrusted comment: minisign public key: E106E938CD763C05
RWQFPHbNOOkG4bXw9P9+wzRhQLwNcBZdgn94TCJyaY7e7CyBYzXXXktB`
)

// Check if there is a new version of the app available and return the url to download it, empty string if there is no update
func (a *App) CheckForUpdates() bool {
	// Check if the latest version is greater than the current version
	return isGreaterVersion(getLatestVersion(), VERSION)
}

// Updates the executable to the one at the given URL
func (a *App) DoUpdate() error {
	// TODO: Binary patching && Checksum verification

	// Get the URL of the new binary for this system
	binaryUrl, err := getBinaryUrl(a.ctx)
	if err != nil {
		return err
	}
	// Get the URL of the signature for this system
	signatureUrl, err := getSignatureUrl(a.ctx)
	if err != nil {
		return err
	}

	fmt.Println("Binary: ", binaryUrl)
	fmt.Println("Signature: ", signatureUrl)

	// Load the binary from the URL
	resp, err := http.Get(binaryUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the body into a byte array to verify its signature
	binary := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return err
		}
		binary = append(binary, buf[:n]...)
	}

	// Initialize the verifier
	verifier := selfupdate.NewVerifier()

	// Load the signature from the URL
	err = verifier.LoadFromURL(
		signatureUrl, // URL of the signature
		PUBLIC_KEY,   // Public key
		http.DefaultTransport,
	)
	if err != nil {
		return err
	}

	// Verify the binary file using the signature
	err = verifier.Verify(binary)
	// If the signature is not valid, return an error
	if err != nil {
		return err
	}

	// Create a reader for the binary (needed for selfupdate.Apply)
	binaryReader := strings.NewReader(string(binary))

	// If the signature is valid, apply the update
	err = selfupdate.Apply(binaryReader, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return rerr
		}
	}
	return err
}

// Get the URL of the binary for the current system
func getBinaryUrl(ctx context.Context) (string, error) {
	// Get the latest release
	release, err := getLatestRelease()
	if err != nil {
		return "", err
	}
	// Get the system info
	system := wailsRuntime.Environment(ctx)
	// Get the right asset for the system
	var systemString string
	if system.Platform == "windows" {
		systemString = "windows"
	} else if system.Platform == "linux" {
		systemString = "linux"
	} else if system.Platform == "darwin" {
		systemString = "macos"
	} else {
		return "", fmt.Errorf("unsupported platform %s", system.Platform)
	}

	// Return the URL of the right asset for this system
	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, systemString) && !strings.HasSuffix(*asset.Name, ".minisig") {
			return *asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no asset found for system %s", systemString)
}

// Get the URL of the signature for the current system
func getSignatureUrl(ctx context.Context) (string, error) {
	// Get the latest release
	release, err := getLatestRelease()
	if err != nil {
		return "", err
	}

	// Get the system info
	system := wailsRuntime.Environment(ctx)
	// Get the right asset for the system
	var systemString string
	if system.Platform == "windows" {
		systemString = "windows"
	} else if system.Platform == "linux" {
		systemString = "linux"
	} else if system.Platform == "darwin" {
		systemString = "macos"
	} else {
		return "", fmt.Errorf("unsupported platform %s", system.Platform)
	}

	// Return the URL of the right asset for this system
	for _, asset := range release.Assets {
		if strings.Contains(*asset.Name, systemString) && strings.HasSuffix(*asset.Name, ".minisig") {
			return *asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no asset found for system %s", systemString)
}

// Checks two versions and returns true if the first one is greater than the second (x.x.x)
func isGreaterVersion(v1, v2 string) bool {
	// Split the versions into their components
	var v1Major, v1Minor, v1Patch int
	var v2Major, v2Minor, v2Patch int
	fmt.Sscanf(v1, "%d.%d.%d", &v1Major, &v1Minor, &v1Patch)
	fmt.Sscanf(v2, "%d.%d.%d", &v2Major, &v2Minor, &v2Patch)

	// Compare the components
	if v1Major > v2Major {
		return true
	}
	if v1Major == v2Major {
		if v1Minor > v2Minor {
			return true
		}
		if v1Minor == v2Minor && v1Patch > v2Patch {
			return true
		}
	}
	return false
}

// Returns the latest version of the app
func getLatestVersion() string {
	// Get the latest release
	release, err := getLatestRelease()
	if err != nil {
		return VERSION
	}
	// Get the tag name
	return release.GetTagName()
}

// Get the release url for the latest version
func (a *App) GetReleaseUrl() (string, error) {
	release, err := getLatestRelease()
	if err != nil {
		return "", err
	}
	return release.GetHTMLURL(), nil
}

// Get the latest release from the GitHub repository
func getLatestRelease() (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "epos-eu", "opensource-desktop")
	if err != nil {
		return nil, err
	}
	return release, nil
}
