package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	permDir os.FileMode = 0o755
)

// Config holds the application configuration.
type Config struct {
	Version string `json:"version"` // Hardcoded version, not from TOML
	Foobar  string `json:"foobar"`  // Value from TOML configuration
}

var cfg *Config

// Init initializes the configuration system.
func Init() error {
	// Set up Viper
	viper.SetConfigName("shed")
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("SHED")
	viper.AutomaticEnv()

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create .shed directory if it doesn't exist.
	shedDir := filepath.Join(homeDir, ".shed")
	if err := os.MkdirAll(shedDir, permDir); err != nil {
		return fmt.Errorf("failed to create .shed directory: %w", err)
	}

	// Add config path
	viper.AddConfigPath(shedDir)

	// Set defaults
	viper.SetDefault("foobar", "default_foobar_value")

	// Read config file (it's okay if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Initialize config struct
	cfg = &Config{
		Version: "1.0.0", // Hardcoded version
		Foobar:  viper.GetString("foobar"),
	}

	return nil
}

// Get returns the current configuration.
func Get() *Config {
	if cfg == nil {
		// Initialize with defaults if not already initialized
		cfg = &Config{
			Version: "0.1.0",
			Foobar:  "default_foobar_value",
		}
	}

	return cfg
}

// GetShedDir returns the .shed directory path.
func GetShedDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".shed"), nil
}

// GetDatabasePath returns the database file path.
func GetDatabasePath() (string, error) {
	shedDir, err := GetShedDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(shedDir, "sheddb"), nil
}
