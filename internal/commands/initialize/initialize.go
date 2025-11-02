package initialize

import (
	"context"
	"errors"
	"fmt"
	"os"

	"h3jfc/shed/internal/logger"
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
	l := logger.Get()
	l.Debug("Starting shed initialization process")
	l.Debug("Checking for existing shed configuration")
	l.Debug("Checking SHED_DIR environment variable")

	location, err := getLocation()
	if err != nil {
		return err
	}

	if location != "" {
		if isValid := validate(location); isValid {
			l.Info("Shed is already initialized", "location", location)

			return nil
		} else {
			l.Debug("Shed configuration found but is invalid", "location", location)

			return ErrConfigInvalid
		}
	}

	// 3. If non exist, Based on OS offer locations
	// 	a. Windows: %USERPROFILE%\.shed or %APPDATA%\.shed
	// 	b. MacOS: $HOME/.shed or $HOME/.config/shed or /etc/shed
	// 	c. Linux: $HOME/.shed or $HOME/.config/shed or /etc/shed
	// 	d. OR ask the user to provide a custom location
	userLocations := getPotentialLocations()

	selectedLocation, err := promptUserForLocation(userLocations)
	if err != nil {
		l.Error("Error selecting location", "error", err)

		return ErrLocationSelection
	}

	err = createShedDirectory(selectedLocation)
	if err != nil {
		l.Error("Error creating shed directory", "error", err)

		return ErrDirectoryCreation
	}

	// 5. Create a default config file in the selected location
	// 6. Create a default database file in the selected location

	// 4. Add ShedDirectory to PATH and ask the user to Add to Path based on common shells (bash, zsh, fish, powershell)
	l.Info("Shed initialized successfully", "location", selectedLocation)
	l.Info("Please add the following line to your shell configuration file to include Shed in your PATH",
		"bash/zsh", "export PATH=\"$PATH:"+selectedLocation+"/bin\"",
		"fish", "set -Ux PATH $PATH "+selectedLocation+"/bin",
		"powershell", "$env:Path += \";"+selectedLocation+"\\bin\"",
	) // make tis conditional based on OS and shell detection

	// 7. Add ShedDir to environment variable SHED_DIR
	l.Info("Please add the following line to your shell configuration file to set SHED_DIR environment variable",
		"bash/zsh", "export SHED_DIR=\""+selectedLocation+"\"",
		"fish", "set -Ux SHED_DIR "+selectedLocation,
		"powershell", "$env:SHED_DIR = \""+selectedLocation+"\"",
	) // make tis conditional based on OS and shell detection
	return ErrNotImplemented
}

func getLocation() (string, error) {
	location := os.Getenv("SHED_DIR")
	if location != "" {
		return location, nil
	}

	var potentialLocations []string

	defaultLocations := getDefaultConfigLocations()

	// check which one exists based on priority
	for _, loc := range defaultLocations {
		if _, err := os.Stat(loc); err == nil {
			location = loc
			potentialLocations = append(potentialLocations, loc)
		}
	}

	l := logger.Get()

	switch len(potentialLocations) {
	case 0:
		l.Debug("No existing shed configuration found. Moving to initialization...")
	case 1:
		l.Debug("One existing shed configuration found!")

		location = potentialLocations[0]
	default:
		l.Debug("Multiple existing shed configurations found", "locations", potentialLocations)

		return location, fmt.Errorf("%w: %v", ErrMultipleConfigs, potentialLocations)
	}

	return location, nil
}

func getDefaultConfigLocations() []string {
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
