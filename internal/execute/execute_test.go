package execute

import (
	"runtime"
	"testing"

	"h3jfc/shed/internal/logger"
)

func TestRun_Success(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "Write-Host 'test output'"
	} else {
		command = "echo 'test output'"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error, got: %v", err)
	}
}

func TestRun_StdoutOutput(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "Write-Host 'line1'; Write-Host 'line2'; Write-Host 'line3'"
	} else {
		command = "echo 'line1' && echo 'line2' && echo 'line3'"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error, got: %v", err)
	}
}

func TestRun_StderrOutput(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "[Console]::Error.WriteLine('error message')"
	} else {
		command = "echo 'error message' >&2"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error, got: %v", err)
	}
}

func TestRun_CommandFailure(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "nonexistentcommand12345"
	} else {
		command = "nonexistentcommand12345"
	}

	err := Run(command)
	if err == nil {
		t.Error("Run() expected error for non-existent command, got nil")
	}
}

func TestRun_ExitCode(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "exit 1"
	} else {
		command = "exit 1"
	}

	err := Run(command)
	if err == nil {
		t.Error("Run() expected error for non-zero exit code, got nil")
	}
}

func TestRun_MixedOutput(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "Write-Host 'stdout message'; [Console]::Error.WriteLine('stderr message')"
	} else {
		command = "echo 'stdout message' && echo 'stderr message' >&2"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error, got: %v", err)
	}
}

func TestRun_ComplexCommand(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		// Test pipe in PowerShell
		command = "Write-Host 'hello' | ForEach-Object { $_.ToUpper() }"
	} else {
		// Test pipe in bash/zsh
		command = "echo 'hello' | tr '[:lower:]' '[:upper:]'"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error for pipe command, got: %v", err)
	}
}

func TestRun_EnvironmentVariables(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv
	// nolint:paralleltest

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	// Set test environment variable
	testValue := "test_value_12345"
	t.Setenv("SHED_TEST_VAR", testValue)

	var command string
	if runtime.GOOS == "windows" {
		command = "Write-Host $env:SHED_TEST_VAR"
	} else {
		command = "echo $SHED_TEST_VAR"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error for env var expansion, got: %v", err)
	}
}

func TestRun_MultilineCommand(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "Write-Host 'line1'; Write-Host 'line2'"
	} else {
		command = "echo 'line1'\necho 'line2'"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error for multiline command, got: %v", err)
	}
}

func TestRun_ShellBuiltins(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		// Get-Location is a PowerShell builtin
		command = "Get-Location"
	} else {
		// pwd is a shell builtin
		command = "pwd"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error for shell builtin, got: %v", err)
	}
}

func TestRun_LongRunningCommand(t *testing.T) {
	t.Parallel()

	// Initialize logger for testing
	logger.New(logger.ModeFromString("message-level"))

	var command string
	if runtime.GOOS == "windows" {
		command = "Start-Sleep -Milliseconds 100; Write-Host 'done'"
	} else {
		command = "sleep 0.1 && echo 'done'"
	}

	err := Run(command)
	if err != nil {
		t.Errorf("Run() expected no error for long-running command, got: %v", err)
	}
}
