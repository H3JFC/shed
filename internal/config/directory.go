package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/h3jfc/shed/internal/logger"
	"github.com/spf13/viper"
)

var (
	ErrConfigInvalid   = errors.New("shed configuration is invalid")
	ErrMultipleConfigs = errors.New("multiple shed configurations found")
	ErrNoPathFound     = errors.New("no shed configuration path found")
)

const (
	defaultConfigName = "config.toml"
)

// GetShedDir returns the .shed directory path.
func GetShedDir() (string, error) {
	// TODO fix
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".shed"), nil
}

// GetDatabasePath returns the database file path.
func GetDatabasePath() (string, error) {
	// TODO fix
	shedDir, err := GetShedDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(shedDir, "sheddb"), nil
}

// FindDir returns the path to the shed configuration.
// it throws an error if it is not found or if it is invalid.
func FindDir() (string, error) {
	p, err := findDirWrap(shedDir, "shedDir")
	if err == nil {
		return p, nil
	}

	return findDirWrap(findDefaultDir, "findDefaultDir")
}

func shedDir() string {
	p := os.Getenv("SHED_DIR")
	if p == "" {
		logger.Debug("SHED_DIR environment variable not set")
	} else {
		logger.Debug("Found SHED_DIR environment variable", "location", p)
	}

	return p
}

func findDefaultDir() string {
	// check which one exists based on priority
	for _, p := range DefaultConfigPaths {
		if _, err := os.Stat(p); err == nil {
			logger.Info("Existing shed path found", "location", p)

			return p
		}
	}

	logger.Debug("No existing shed path found.")

	return ""
}

func findDirWrap(fp func() string, function string) (string, error) {
	if p := fp(); p != "" {
		if isValid := validatePath(p); !isValid {
			logger.Debug("Shed configuration found but is invalid", "location", p, "function", function)

			return "", fmt.Errorf("%w for function %s", ErrConfigInvalid, function)
		}

		logger.Debug("Shed is initialized", "location", p, "function", function)

		return p, nil
	}

	logger.Debug("No shed configuration found in standard locations", "function", function)

	return "", fmt.Errorf("%w for function %s", ErrNoPathFound, function)
}

// validatePath checks if a path contains a valid shed directory structure.
func validatePath(p string) bool {
	// Check if the provided path is a directory
	info, err := os.Stat(p)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check if SQLite database file exists
	dbPath := filepath.Join(p, defaultDBName)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}

	// Check if config.toml exists
	configPath := filepath.Join(p, defaultConfigName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	// Validate config.toml structure using viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return false
	}

	// Validate required [shed-db] section and password field
	if !v.IsSet("shed-db.password") {
		return false
	}

	// Validate that password is not empty
	password := v.GetString("shed-db.password")

	return password != ""
}
