package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const versionFile = "version.json" // Relative path to version.json

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}

	// Get the absolute path to version.json
	versionFilePath := filepath.Join(wd, versionFile)

	// Load the current version from the file
	version, err := loadVersion(versionFilePath)
	if err != nil {
		fmt.Println("Error loading version:", err)
		os.Exit(1)
	}

	// Parse command-line arguments
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Error: expects only 1 cmdline argument; go bump_version.go [major|minor|patch]")
		os.Exit(1)
	}

	arg := strings.ToLower(args[1])
	if arg != "major" && arg != "minor" && arg != "patch" {
		fmt.Println("Error: argument must be one of [major|minor|patch]; go bump_version.go [major|minor|patch]")
		os.Exit(1)
	}

	// Bump version based on the argument
	switch arg {
	case "patch":
		version.Patch++
	case "minor":
		version.Minor++
		version.Patch = 0
	case "major":
		version.Major++
		version.Minor = 0
		version.Patch = 0
	}

	// Format the version string
	newVersion := fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)

	// Save the updated version to the file
	err = saveVersion(versionFilePath, version)
	if err != nil {
		fmt.Println("Error saving version:", err)
		os.Exit(1)
	}

	// Add and commit version.json before creating the Git tag
	err = addAndCommitVersion(versionFilePath)
	if err != nil {
		fmt.Println("Error committing version changes:", err)
		os.Exit(1)
	}

	// Create the Git tag
	err = createGitTag(newVersion)
	if err != nil {
		fmt.Printf("Error creating Git tag: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Version bumped to %s and Git tag created.\n", newVersion)
}

// loadVersion loads the current version from the version file
func loadVersion(versionFilePath string) (Version, error) {
	var version Version

	// Check if the version file exists
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		// If the file doesn't exist, start with version 0.0.1
		version = Version{Major: 0, Minor: 0, Patch: 1}
		return version, nil
	}

	// Open the file
	file, err := os.Open(versionFilePath)
	if err != nil {
		return version, fmt.Errorf("failed to open version file: %w", err)
	}
	defer file.Close()

	// Read the file content
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&version)
	if err != nil {
		return version, fmt.Errorf("failed to decode version file: %w", err)
	}

	return version, nil
}

// saveVersion saves the updated version to the version file
func saveVersion(versionFilePath string, version Version) error {
	file, err := os.Create(versionFilePath)
	if err != nil {
		return fmt.Errorf("failed to create version file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(version)
	if err != nil {
		return fmt.Errorf("failed to encode version data: %w", err)
	}

	return nil
}

// Add and commit changes
func addAndCommitVersion(versionFilePath string) error {
	// Run git add command to stage the version file
	cmd := exec.Command("git", "add", versionFilePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to stage version file: %w", err)
	}

	// Run git commit command to commit the changes
	cmd = exec.Command("git", "commit", "-m", "Bump version")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

// createGitTag creates a Git tag with the specified version
func createGitTag(version string) error {
	// Run git tag command to create a new tag
	cmd := exec.Command("git", "tag", version)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}
	