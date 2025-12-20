//go:build linux

package execute

import (
	"os"
	"path/filepath"
	"strings"
)

// detectShellPlatform detects the appropriate shell for Linux.
func detectShellPlatform() ShellConfig {
	// 1. Check $SHELL environment variable
	if shellPath := os.Getenv("SHELL"); shellPath != "" {
		if config := shellConfigFromPath(shellPath); config.Path != "" {
			return config
		}
	}

	// 2. Check for shell config files
	shells := detectShellsByConfigFiles()
	if len(shells) > 0 {
		return shellConfigFromName(shells[0])
	}

	// 3. Default to bash (most common on Linux)
	return ShellConfig{
		Name: "bash",
		Path: "/bin/bash",
		Args: []string{"-c"},
	}
}

// detectShellsByConfigFiles checks which shell config files exist in the user's home directory.
func detectShellsByConfigFiles() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
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

	// Check in preferred order for Linux: bash, zsh, fish
	for _, shell := range []string{"bash", "zsh", "fish"} {
		paths := shellConfigs[shell]
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				detected = append(detected, shell)
				break // found one config for this shell, move to next
			}
		}
	}

	return detected
}

// shellConfigFromPath creates a ShellConfig from a full shell path.
func shellConfigFromPath(shellPath string) ShellConfig {
	// Extract shell name from path (e.g., "/bin/bash" -> "bash")
	shellName := filepath.Base(shellPath)

	// Validate that the shell exists and is executable
	if info, err := os.Stat(shellPath); err != nil || info.IsDir() {
		return ShellConfig{} // Invalid path
	}

	return shellConfigFromName(shellName)
}

// shellConfigFromName creates a ShellConfig for a known shell name.
func shellConfigFromName(name string) ShellConfig {
	// Remove potential version suffixes (e.g., "bash5" -> "bash")
	baseName := strings.TrimRight(name, "0123456789")

	switch baseName {
	case "bash":
		return ShellConfig{
			Name: "bash",
			Path: "/bin/bash",
			Args: []string{"-c"},
		}
	case "zsh":
		return ShellConfig{
			Name: "zsh",
			Path: "/bin/zsh",
			Args: []string{"-c"},
		}
	case "fish":
		return ShellConfig{
			Name: "fish",
			Path: "/usr/bin/fish",
			Args: []string{"-c"},
		}
	default:
		// For unknown shells, assume POSIX-compatible
		return ShellConfig{
			Name: name,
			Path: "/bin/" + name,
			Args: []string{"-c"},
		}
	}
}
