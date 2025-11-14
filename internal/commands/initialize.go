package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"h3jfc/shed/internal/config"
	"h3jfc/shed/internal/logger"
)

const retryAttempts = 3

var (
	ErrNotImplemented     = errors.New("not implemented yet")
	ErrConfigInvalid      = errors.New("shed configuration is invalid")
	ErrLocationSelection  = errors.New("error selecting shed location")
	ErrDirectoryCreation  = errors.New("error creating shed directory")
	ErrMultipleConfigs    = errors.New("multiple shed configurations found")
	ErrMaxAttemptsReached = errors.New("maximum attempts reached for location selection")
	ErrInvalidInput       = errors.New("invalid input: please enter a number")
	ErrInvalidChoice      = errors.New("invalid choice")
)

func Init(_ context.Context) error {
	dir, err := promptUserDirWithRetry(config.DefaultConfigPaths, retryAttempts)
	if err != nil {
		logger.Error("Error selecting location", "error", err)

		return ErrLocationSelection
	}

	if err := config.CreateShedDirectory(dir); err != nil {
		logger.Error("Error creating shed directory and db", "error", err)
		os.RemoveAll(dir) // cleanup on failure

		return ErrDirectoryCreation
	}

	logShellInstructions(dir)

	return nil
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
		return "", ErrInvalidInput
	}

	// Validate the choice
	if choice < 1 || choice > len(locations) {
		return "", fmt.Errorf("%w: please select a number between 1 and %d", ErrInvalidChoice, len(locations))
	}

	return locations[choice-1], nil
}

func promptUserDirWithRetry(locations []string, maxAttempts int) (string, error) {
	for attempt := range maxAttempts {
		location, err := promptUserDir(locations)
		if err == nil {
			return location, nil
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		if attempt < maxAttempts-1 {
			fmt.Println("Please try again.")
		}
	}

	return "", ErrMaxAttemptsReached
}

func logShellInstructions(dir string) {
	logger.Info("Shed initialized successfully", "location", dir)

	shells := detectShellsByConfigFiles()

	if len(shells) == 1 {
		// Only one shell detected - provide specific instructions
		shell := shells[0]
		shedDirInstr, shedDirConfig := getShedDirInstruction(dir, shell)

		logger.Info(fmt.Sprintf("Detected %s - set SHED_DIR:", shell),
			"instruction", shedDirInstr,
			"config_file", shedDirConfig,
		)
	} else {
		// Multiple shells detected - provide all instructions
		logger.Info(fmt.Sprintf("Detected shells: %v - set SHED_DIR:", shells))

		for _, shell := range shells {
			shedDirInstr, shedDirConfig := getShedDirInstruction(dir, shell)
			logger.Info("  "+shell,
				"instruction", shedDirInstr,
				"config_file", shedDirConfig,
			)
		}
	}
}

// detectShellsByConfigFiles checks which shell config files exist.
func detectShellsByConfigFiles() []string {
	if runtime.GOOS == "windows" {
		return []string{"powershell"}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{"bash"} // fallback
	}

	shellConfigs := map[string][]string{
		"bash": {
			filepath.Join(homeDir, ".bashrc"),
			filepath.Join(homeDir, ".bash_profile"),
			filepath.Join(homeDir, ".profile"),
		},
		"zsh": {
			filepath.Join(homeDir, ".zshrc"),
			filepath.Join(homeDir, ".zshenv"),
		},
		"fish": {
			filepath.Join(homeDir, ".config", "fish", "config.fish"),
		},
	}

	var detected []string

	for shell, paths := range shellConfigs {
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				detected = append(detected, shell)

				break // found one config for this shell, move to next
			}
		}
	}

	// If nothing detected, default to bash
	if len(detected) == 0 {
		detected = []string{"bash"}
	}

	return detected
}

// getShedDirInstruction returns the appropriate SHED_DIR instruction for a shell.
func getShedDirInstruction(dir, shell string) (instruction, configFile string) {
	switch shell {
	case "fish":
		instruction = "set -Ux SHED_DIR " + dir
		configFile = "~/.config/fish/config.fish"
	case "zsh":
		instruction = fmt.Sprintf("export SHED_DIR=\"%s\"", dir)
		configFile = "~/.zshrc"
	case "powershell":
		instruction = fmt.Sprintf("$env:SHED_DIR = \"%s\"", filepath.ToSlash(dir))
		configFile = "$PROFILE"
	case "bash":
		fallthrough
	default:
		instruction = fmt.Sprintf("export SHED_DIR=\"%s\"", dir)
		configFile = "~/.bashrc"
	}

	return instruction, configFile
}
