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

const version = "0.0.1"

// Check if there is a new version of the app available and return the url to download it, empty string if there is no update
func (a *App) CheckForUpdates(system wailsRuntime.EnvironmentInfo) (string, error) {
	//TODO: Check for updates
	return getLatestGithubReleaseUrl(system)
}

// Updates the executable to the one at the given URL
func (a *App) DoUpdate(url string) error {
	// TODO: Binary patching && Checksum verification
	// Request the new file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return rerr
		}
	}
	return err
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
	HtmlUrl string  `json:"html_url"`
}

// Get the latest release url for the app
func getLatestGithubReleaseUrl(system wailsRuntime.EnvironmentInfo) (string, error) {
	// Get the latest release
	client := github.NewClient(nil)
	// TODO: use the correct owner and repo
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "epos-eu", "opensource-docker")
	if err != nil {
		return "", err
	}
	// Get the tag name
	tagName := release.GetTagName()

	// Check if the latest release is newer than the current version
	if isGreaterVersion(tagName, version) {
		// If on mac, return the html url for the manual download
		if system.Platform == "darwin" {
			return release.GetHTMLURL(), nil
		}
		systemString := system.Platform + "_" + system.Arch
		// Return the URL of the right asset for this system
		for _, asset := range release.Assets {
			if strings.Contains(*asset.Name, systemString) {
				return *asset.BrowserDownloadURL, nil
			}
		}
		return "", fmt.Errorf("no asset found for system %s", systemString)
	}

	// No update available
	return "", nil
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
