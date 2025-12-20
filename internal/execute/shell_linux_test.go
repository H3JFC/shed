//go:build linux

package execute

import (
	"os"
	"testing"
)

func TestDetectShellPlatform_Linux(t *testing.T) {
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

	// On Linux, should default to bash or detect from $SHELL
	if config.Name != "bash" && config.Name != "zsh" && config.Name != "fish" {
		t.Logf("Unexpected shell detected: %s (this is OK if it's a custom shell)", config.Name)
	}
}

func TestShellConfigFromName_Bash(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName("bash")

	if config.Name != "bash" {
		t.Errorf("shellConfigFromName(bash) name = %v, want bash", config.Name)
	}

	if config.Path != "/bin/bash" {
		t.Errorf("shellConfigFromName(bash) path = %v, want /bin/bash", config.Path)
	}

	if len(config.Args) != 1 || config.Args[0] != "-c" {
		t.Errorf("shellConfigFromName(bash) args = %v, want [-c]", config.Args)
	}
}

func TestShellConfigFromName_Zsh(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName("zsh")

	if config.Name != "zsh" {
		t.Errorf("shellConfigFromName(zsh) name = %v, want zsh", config.Name)
	}

	if config.Path != "/bin/zsh" {
		t.Errorf("shellConfigFromName(zsh) path = %v, want /bin/zsh", config.Path)
	}

	if len(config.Args) != 1 || config.Args[0] != "-c" {
		t.Errorf("shellConfigFromName(zsh) args = %v, want [-c]", config.Args)
	}
}

func TestShellConfigFromName_Fish(t *testing.T) {
	t.Parallel()

	config := shellConfigFromName("fish")

	if config.Name != "fish" {
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

	// Test with /bin/bash which should exist on Linux
	config := shellConfigFromPath("/bin/bash")

	if config.Name != "bash" {
		t.Errorf("shellConfigFromPath(/bin/bash) name = %v, want bash", config.Name)
	}

	if config.Path != "/bin/bash" {
		t.Errorf("shellConfigFromPath(/bin/bash) path = %v, want /bin/bash", config.Path)
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

func TestDetectShellsByConfigFiles_Linux(t *testing.T) {
	// This test checks actual files, so don't run in parallel
	// nolint:paralleltest

	shells := detectShellsByConfigFiles()

	// Should detect at least one shell (bash should have config files on most Linux systems)
	if len(shells) == 0 {
		t.Log("No shell config files detected (this is OK in CI environments)")
	}

	// Verify detected shells are valid
	for _, shell := range shells {
		if shell != "bash" && shell != "zsh" && shell != "fish" {
			t.Errorf("detectShellsByConfigFiles() returned unexpected shell: %s", shell)
		}
	}
}

func TestGetShellConfig_Caching(t *testing.T) {
	// nolint:paralleltest
	// This test manipulates global state, so don't run in parallel

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

func TestSetShellConfig_CustomShell(t *testing.T) {
	// nolint:paralleltest
	// This test manipulates global state, so don't run in parallel

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

func TestDetectShellPlatform_WithSHELLEnv(t *testing.T) {
	// nolint:paralleltest
	// This test manipulates environment variables

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

	// Set $SHELL to /bin/zsh
	os.Setenv("SHELL", "/bin/zsh")

	config := detectShellPlatform()

	if config.Name != "zsh" {
		t.Errorf("detectShellPlatform() with SHELL=/bin/zsh got name = %v, want zsh", config.Name)
	}
}
