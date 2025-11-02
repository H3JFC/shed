package initialize

import (
	"context"
	"errors"

	"h3jfc/shed/internal/config"
	"h3jfc/shed/internal/logger"
	libos "h3jfc/shed/lib/os"
)

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

	p, err := config.FindPath()
	if err != nil && !errors.Is(err, config.ErrNoPathFound) {
		return err
	}

	if err == nil && p != "" {
		logger.Info("Shed is already initialized at location", "location", p)
		logger.Debug("Initialization aborted to prevent overwriting existing configuration")

		return nil
	}

	// 3. If non exist, Based on OS offer locations
	// 	a. Windows: %USERPROFILE%\.shed or %APPDATA%\.shed
	// 	b. MacOS: $HOME/.shed or $HOME/.config/shed or /etc/shed
	// 	c. Linux: $HOME/.shed or $HOME/.config/shed or /etc/shed
	// 	d. OR ask the user to provide a custom location
	userLocations := getPotentialLocations()

	selectedLocation, err := promptUserForLocation(userLocations)
	if err != nil {
		logger.Error("Error selecting location", "error", err)

		return ErrLocationSelection
	}

	err = createShedDirectory(selectedLocation)
	if err != nil {
		logger.Error("Error creating shed directory", "error", err)

		return ErrDirectoryCreation
	}

	// 5. Create a default config file in the selected location
	// 6. Create a default database file in the selected location

	// 4. Add ShedDirectory to PATH and ask the user to Add to Path based on common shells (bash, zsh, fish, powershell)
	logger.Info("Shed initialized successfully", "location", selectedLocation)
	logger.Info("Please add the following line to your shell configuration file to include Shed in your PATH",
		"bash/zsh", "export PATH=\"$PATH:"+selectedLocation+"/bin\"",
		"fish", "set -Ux PATH $PATH "+selectedLocation+"/bin",
		"powershell", "$env:Path += \";"+selectedLocation+"\\bin\"",
	) // make tis conditional based on OS and shell detection

	// 7. Add ShedDir to environment variable SHED_DIR
	logger.Info("Please add the following line to your shell configuration file to set SHED_DIR environment variable",
		"bash/zsh", "export SHED_DIR=\""+selectedLocation+"\"",
		"fish", "set -Ux SHED_DIR "+selectedLocation,
		"powershell", "$env:SHED_DIR = \""+selectedLocation+"\"",
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
	panic("implement me")
	return []string{}
}

func promptUserForLocation(locations []string) (string, error) {
	panic("implement me")
	return "", nil
}

func createShedDirectory(location string) error {
	panic("implement me")
	return nil
}
