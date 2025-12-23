package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/h3jfc/shed/lib/sqlite3"
	"golang.org/x/term"
)

const (
	defaultDBName    = "shed.db"
	defaultDirPerms  = 0o755
	defaultFilePerms = 0o644
)

var (
	ErrDirectoryCreation = errors.New("failed to create shed directory")
	ErrPasswordMismatch  = errors.New("passwords do not match")
	ErrEmptyPassword     = errors.New("password cannot be empty")
	ErrNoLocations       = errors.New("no locations provided")
	ErrInvalidInput      = errors.New("invalid input: please enter a number")
	ErrInvalidChoice     = errors.New("invalid choice")
)

// CreateShedDirectory creates the shed directory structure and initializes required files.
func CreateShedDirectory(path string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(path, defaultDirPerms); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	// Prompt for database password
	password, err := promptForPassword()
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Create config.toml with the password
	if err := createConfigFile(path, password); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Create empty database file (will be initialized later with encryption)
	dbPath := filepath.Join(path, defaultDBName)
	if err := sqlite3.MigrateShedDB(dbPath, password); err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}

	return nil
}

// promptForPassword prompts the user to enter and confirm a password.
func promptForPassword() (string, error) {
	fmt.Print("Enter database password: ")

	password, err := readPassword()
	if err != nil {
		return "", err
	}

	fmt.Println() // Add newline after password input

	if strings.TrimSpace(password) == "" {
		return "", ErrEmptyPassword
	}

	fmt.Print("Confirm database password: ")

	confirmPassword, err := readPassword()
	if err != nil {
		return "", err
	}

	fmt.Println() // Add newline after password input

	if password != confirmPassword {
		return "", ErrPasswordMismatch
	}

	return password, nil
}

// readPassword reads a password from stdin without echoing.
func readPassword() (string, error) {
	// nolint:unconvert // Required for Windows compatibility where syscall.Stdin is uintptr
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	return string(bytePassword), nil
}

// createConfigFile creates a config.toml file with the database password.
func createConfigFile(dirPath, password string) error {
	configPath := filepath.Join(dirPath, defaultConfigName)
	dbPath := filepath.Join(dirPath, defaultDBName)

	// Use forward slashes for cross-platform compatibility in TOML
	// (backslashes are escape characters in TOML double-quoted strings)
	dbPathNormalized := filepath.ToSlash(dbPath)

	configContent := fmt.Sprintf(`[shed-db]
password = "%s"
location = "%s"

[settings]
# Add other configuration settings here
`, password, dbPathNormalized)

	if err := os.WriteFile(configPath, []byte(configContent), defaultFilePerms); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// createEmptyFile creates an empty file at the specified path.
func createEmptyFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer f.Close()

	return nil
}

// promptUserForLocation prompts the user to select from a list of locations.
func promptUserForLocation(locations []string) (string, error) {
	if len(locations) == 0 {
		return "", ErrNoLocations
	}

	// Display the list of locations
	fmt.Println("Please select a configuration location:")

	for i, location := range locations {
		fmt.Printf("%d) %s\n", i+1, location)
	}

	fmt.Print("\nEnter your choice (number): ")

	// Read user input
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	// Parse the input
	input = strings.TrimSpace(input)

	var choice int
	if _, err := fmt.Sscanf(input, "%d", &choice); err != nil {
		return "", ErrInvalidInput
	}

	// Validate the choice
	if choice < 1 || choice > len(locations) {
		return "", fmt.Errorf("%w: please select a number between 1 and %d", ErrInvalidChoice, len(locations))
	}

	return locations[choice-1], nil
}
