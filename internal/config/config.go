package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"h3jfc/shed/internal/logger"
)

var (
	ErrConfigInvalid   = errors.New("shed configuration is invalid")
	ErrMultipleConfigs = errors.New("multiple shed configurations found")
	ErrNoPathFound     = errors.New("no shed configuration path found")
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
			logger.Debug("Existing shed path found", "location", p)

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

func validatePath(p string) bool {
	// Check if the provided path is a directory
	info, err := os.Stat(p)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check if SQLite database file exists
	dbPath := filepath.Join(p, "shed.db") // or whatever your db filename is
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}

	// Check if config.toml exists
	configPath := filepath.Join(p, "config.toml")
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

	// Validate required [shed_db] section and password field
	if !v.IsSet("shed_db.password") {
		return false
	}

	// Optional: Validate that password is not empty
	password := v.GetString("shed_db.password")

	return password != ""
}

func Create(location string) error {
	panic("implement me")
	return nil
}
