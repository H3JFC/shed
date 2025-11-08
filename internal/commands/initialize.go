package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"h3jfc/shed/internal/config"
	"h3jfc/shed/internal/logger"
	libos "h3jfc/shed/lib/os"
)

const retryAttempts = 3

var (
	ErrNotImplemented    = errors.New("not implemented yet")
	ErrConfigInvalid     = errors.New("shed configuration is invalid")
	ErrLocationSelection = errors.New("error selecting shed location")
	ErrDirectoryCreation = errors.New("error creating shed directory")
	ErrMultipleConfigs   = errors.New("multiple shed configurations found")
)

// TODO init SQLITE database.
func Init(_ context.Context) error {
	logger.Debug("Starting shed initialization process")
	logger.Debug("Checking for existing shed configuration")
	logger.Debug("Checking SHED_DIR environment variable")

	p, err := config.FindDir()
	if err != nil && !errors.Is(err, config.ErrNoPathFound) {
		return err
	}

	if err == nil && p != "" {
		logger.Info("Shed is already initialized at location", "location", p)
		logger.Debug("Initialization aborted to prevent overwriting existing configuration")

		return nil
	}

	dir, err := promptUserDirWithRetry(config.DefaultConfigPaths, retryAttempts)
	if err != nil {
		logger.Error("Error selecting location", "error", err)

		return ErrLocationSelection
	}

	err = config.Create(dir)
	if err != nil {
		logger.Error("Error creating shed directory", "error", err)

		return ErrDirectoryCreation
	}

	// 5. Create a default config file in the selected location
	// 6. Create a default database file in the selected location

	// 4. Add ShedDirectory to PATH and ask the user to Add to Path based on common shells (bash, zsh, fish, powershell)
	logger.Info("Shed initialized successfully", "location", dir)
	logger.Info("Please add the following line to your shell configuration file to include Shed in your PATH",
		"bash/zsh", "export PATH=\"$PATH:"+dir+"/bin\"",
		"fish", "set -Ux PATH $PATH "+dir+"/bin",
		"powershell", "$env:Path += \";"+dir+"\\bin\"",
	) // make tis conditional based on OS and shell detection

	// 7. Add ShedDir to environment variable SHED_DIR
	logger.Info("Please add the following line to your shell configuration file to set SHED_DIR environment variable",
		"bash/zsh", "export SHED_DIR=\""+dir+"\"",
		"fish", "set -Ux SHED_DIR "+dir,
		"powershell", "$env:SHED_DIR = \""+dir+"\"",
	) // make tis conditional based on OS and shell detection
	return ErrNotImplemented
}

func defaultConfigLocations(o libos.OS) []string {
	panic("implement me")
	return []string{"~/.shed"}
}

func validate(loc string) bool {
	panic("implement me")
	return false
}

func getPotentialLocations() []string {
	return config.DefaultConfigPaths
}

func promptUserDir(locations []string) (string, error) {
	if len(locations) == 0 {
		logger.Error("No configuration locations provided to promptUserForLocation")
		// invariant violation. Each os should have at least one default location
		panic("Well this should never happen")
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
	choice, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("invalid input: please enter a number")
	}

	// Validate the choice
	if choice < 1 || choice > len(locations) {
		return "", fmt.Errorf("invalid choice: please select a number between 1 and %d", len(locations))
	}

	return locations[choice-1], nil
}

func promptUserDirWithRetry(locations []string, maxAttempts int) (string, error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		location, err := promptUserDir(locations)
		if err == nil {
			return location, nil
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if attempt < maxAttempts-1 {
			fmt.Println("Please try again.")
		}
	}
	return "", fmt.Errorf("max attempts reached")
}
