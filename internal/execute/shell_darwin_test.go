//go:build darwin

package execute

import (
	"os"
	"testing"
)

const (
	zshShell  = "zsh"
	bashShell = "bash"
	fishShell = "fish"
)

func TestDetectShellPlatform_Darwin(t *testing.T) {
	t.Parallel()

	// Reset shell config before test
	ResetShellConfig()
	defer ResetShellConfig()

	config := detectShellPlatform()

	// Should detect a shell
	if config.Name == "" {
		t.Error("detectShellPlatform() returned empty shell name")
	}

	if config.Path == "" {
		t.Error("detectShellPlatform() returned empty shell path")
	}

	if len(config.Args) == 0 {
		t.Error("detectShellPlatform() returned empty shell args")
	}

	// On macOS, should default to zsh or detect from $SHELL
	if config.Name != zshShell && config.Name != bashShell && config.Name != fishShell {
		t.Logf("Unexpected shell detected: %s (this is OK if it's a custom shell)", config.Name)
	}
}

func TestShellConfigFromName_Zsh(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName(zshShell)

	if config.Name != zshShell {
		t.Errorf("shellConfigFromName(zsh) name = %v, want zsh", config.Name)
	}

	if config.Path != "/bin/zsh" {
		t.Errorf("shellConfigFromName(zsh) path = %v, want /bin/zsh", config.Path)
	}

	if len(config.Args) != 1 || config.Args[0] != "-c" {
		t.Errorf("shellConfigFromName(zsh) args = %v, want [-c]", config.Args)
	}
}

func TestShellConfigFromName_Bash(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName(bashShell)

	if config.Name != bashShell {
		t.Errorf("shellConfigFromName(bash) name = %v, want bash", config.Name)
	}

	if config.Path != "/bin/bash" {
		t.Errorf("shellConfigFromName(bash) path = %v, want /bin/bash", config.Path)
	}

	if len(config.Args) != 1 || config.Args[0] != "-c" {
		t.Errorf("shellConfigFromName(bash) args = %v, want [-c]", config.Args)
	}
}

func TestShellConfigFromName_Fish(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName(fishShell)

	if config.Name != fishShell {
		t.Errorf("shellConfigFromName(fish) name = %v, want fish", config.Name)
	}

	if config.Path != "/usr/bin/fish" {
		t.Errorf("shellConfigFromName(fish) path = %v, want /usr/bin/fish", config.Path)
	}

	if len(config.Args) != 1 || config.Args[0] != "-c" {
		t.Errorf("shellConfigFromName(fish) args = %v, want [-c]", config.Args)
	}
}

func TestShellConfigFromPath_ValidPath(t *testing.T) {
	t.Parallel()

	// Test with /bin/zsh which should exist on macOS
	config := shellConfigFromPath("/bin/zsh")

	if config.Name != zshShell {
		t.Errorf("shellConfigFromPath(/bin/zsh) name = %v, want zsh", config.Name)
	}

	if config.Path != "/bin/zsh" {
		t.Errorf("shellConfigFromPath(/bin/zsh) path = %v, want /bin/zsh", config.Path)
	}
}

func TestShellConfigFromPath_InvalidPath(t *testing.T) {
	t.Parallel()

	// Test with non-existent path
	config := shellConfigFromPath("/nonexistent/shell")

	if config.Path != "" {
		t.Errorf("shellConfigFromPath(invalid) should return empty config, got %+v", config)
	}
}

// nolint:paralleltest // This test checks actual files, so don't run in parallel
func TestDetectShellsByConfigFiles_Darwin(t *testing.T) {
	shells := detectShellsByConfigFiles()

	// Should detect at least one shell (zsh or bash should have config files on macOS)
	if len(shells) == 0 {
		t.Log("No shell config files detected (this is OK in CI environments)")
	}

	// Verify detected shells are valid
	for _, shell := range shells {
		if shell != bashShell && shell != zshShell && shell != fishShell {
			t.Errorf("detectShellsByConfigFiles() returned unexpected shell: %s", shell)
		}
	}
}

// nolint:paralleltest // This test manipulates global state, so don't run in parallel
func TestGetShellConfig_Caching(t *testing.T) {
	ResetShellConfig()
	defer ResetShellConfig()

	// First call should detect shell
	config1 := GetShellConfig()

	// Second call should return cached value
	config2 := GetShellConfig()

	if config1.Name != config2.Name {
		t.Errorf("GetShellConfig() not caching properly: first=%s, second=%s", config1.Name, config2.Name)
	}

	if config1.Path != config2.Path {
		t.Errorf("GetShellConfig() not caching properly: first=%s, second=%s", config1.Path, config2.Path)
	}
}

// nolint:paralleltest // This test manipulates global state, so don't run in parallel
func TestSetShellConfig_CustomShell(t *testing.T) {
	ResetShellConfig()
	defer ResetShellConfig()

	customConfig := ShellConfig{
		Name: "custom",
		Path: "/custom/shell",
		Args: []string{"-x"},
	}

	SetShellConfig(customConfig)

	retrieved := GetShellConfig()

	if retrieved.Name != "custom" {
		t.Errorf("SetShellConfig() name = %v, want custom", retrieved.Name)
	}

	if retrieved.Path != "/custom/shell" {
		t.Errorf("SetShellConfig() path = %v, want /custom/shell", retrieved.Path)
	}

	if len(retrieved.Args) != 1 || retrieved.Args[0] != "-x" {
		t.Errorf("SetShellConfig() args = %v, want [-x]", retrieved.Args)
	}
}

// nolint:paralleltest // This test manipulates environment variables
func TestDetectShellPlatform_WithSHELLEnv(t *testing.T) {
	ResetShellConfig()
	defer ResetShellConfig()

	// Save original $SHELL
	originalShell := os.Getenv("SHELL")

	defer func() {
		if originalShell != "" {
			os.Setenv("SHELL", originalShell)
		} else {
			os.Unsetenv("SHELL")
		}
	}()

	// Set $SHELL to /bin/bash
	os.Setenv("SHELL", "/bin/bash")

	config := detectShellPlatform()

	if config.Name != bashShell {
		t.Errorf("detectShellPlatform() with SHELL=/bin/bash got name = %v, want bash", config.Name)
	}
}
