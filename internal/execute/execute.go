// Package execute provides functionality for executing shell commands
// with proper logging of stdout and stderr.
//
// # Shell Detection
//
// The package automatically detects the appropriate shell for the platform:
//   - macOS: Uses $SHELL env var, or detects from config files (.zshrc, .bashrc),
//     defaults to /bin/zsh (macOS default since Catalina)
//   - Linux: Uses $SHELL env var, or detects from config files (.bashrc, .zshrc),
//     defaults to /bin/bash (most common on Linux)
//   - Windows: Prefers PowerShell Core (pwsh), then Windows PowerShell,
//     falls back to cmd.exe
//
// The detected shell is cached after first use for performance.
//
// # Usage
//
// Basic command execution:
//
//	err := execute.Run("echo 'Hello, World!'")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Complex commands with pipes, redirects, and shell features are supported:
//
//	err := execute.Run("cat file.txt | grep 'pattern' > output.txt")
//	err := execute.Run("for i in {1..5}; do echo $i; done")
//
// Shell builtins and environment variable expansion work as expected:
//
//	err := execute.Run("cd /tmp && pwd")
//	err := execute.Run("echo $HOME")
//
// The function blocks until the command completes. Stdout is logged at Info level,
// stderr is logged at Error level.
package execute

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"h3jfc/shed/internal/logger"
)

const (
	numWaitGroups = 2
)

// Run executes a command through the system shell and logs output.
//
// The command is executed in the user's default shell (bash, zsh, PowerShell, etc.)
// which is automatically detected. Stdout is logged at Info level, stderr at Error level.
//
// The function blocks until the command completes. If the command returns a non-zero
// exit code, an error is returned.
//
// Example:
//
//	err := execute.Run("echo 'Hello, World!'")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err := execute.Run("ls -la | grep '.go'")
func Run(command string) error {
	// Get shell configuration (cached after first call)
	shellConfig := GetShellConfig()

	// Create command with proper shell invocation
	// #nosec G204 -- Command execution is the intended functionality of this package
	cmd := exec.Command(shellConfig.Path, append(shellConfig.Args, command)...)

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Use a WaitGroup to wait for both stdout and stderr readers to finish
	var wg sync.WaitGroup

	wg.Add(numWaitGroups)

	// Stream stdout to logger.Info
	go func() {
		defer wg.Done()

		streamToLogger(stdout, logger.Info)
	}()

	// Stream stderr to logger.Error
	go func() {
		defer wg.Done()

		streamToLogger(stderr, logger.Error)
	}()

	// Wait for all output to be read
	wg.Wait()

	// Wait for the command to finish and check for errors
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// streamToLogger reads from an io.Reader line by line and logs each line
// using the provided log function.
func streamToLogger(reader io.Reader, logFunc func(string, ...any)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logFunc(scanner.Text())
	}
}
