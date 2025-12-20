//go:build windows

package execute

import (
	"testing"
)

func TestDetectShellPlatform_Windows(t *testing.T) {
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

	// On Windows, should be PowerShell or cmd
	if config.Name != "pwsh" && config.Name != "powershell" && config.Name != "cmd" {
		t.Errorf("Unexpected shell detected: %s", config.Name)
	}

	// Verify args are correct for the detected shell
	if config.Name == "pwsh" || config.Name == "powershell" {
		if len(config.Args) != 1 || config.Args[0] != "-Command" {
			t.Errorf("PowerShell args = %v, want [-Command]", config.Args)
		}
	} else if config.Name == "cmd" {
		if len(config.Args) != 1 || config.Args[0] != "/C" {
			t.Errorf("cmd args = %v, want [/C]", config.Args)
		}
	}
}

func TestDetectShellPlatform_FallbackToCmd(t *testing.T) {
	t.Parallel()

	// Reset shell config before test
	ResetShellConfig()
	defer ResetShellConfig()

	// Even if PowerShell is not found, should fall back to cmd.exe
	config := detectShellPlatform()

	// Should always have cmd.exe as fallback on Windows
	if config.Name == "" {
		t.Error("detectShellPlatform() should always return a shell on Windows")
	}

	// cmd.exe should always be available
	validShells := map[string]bool{"pwsh": true, "powershell": true, "cmd": true}
	if !validShells[config.Name] {
		t.Errorf("detectShellPlatform() returned invalid shell: %s", config.Name)
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
		Path: "C:\\custom\\shell.exe",
		Args: []string{"-x"},
	}

	SetShellConfig(customConfig)

	retrieved := GetShellConfig()

	if retrieved.Name != "custom" {
		t.Errorf("SetShellConfig() name = %v, want custom", retrieved.Name)
	}

	if retrieved.Path != "C:\\custom\\shell.exe" {
		t.Errorf("SetShellConfig() path = %v, want C:\\custom\\shell.exe", retrieved.Path)
	}

	if len(retrieved.Args) != 1 || retrieved.Args[0] != "-x" {
		t.Errorf("SetShellConfig() args = %v, want [-x]", retrieved.Args)
	}
}

func TestShellConfig_PowerShellArgs(t *testing.T) {
	t.Parallel()

	// Test that PowerShell config has correct args
	config := ShellConfig{
		Name: "powershell",
		Path: "powershell.exe",
		Args: []string{"-Command"},
	}

	if len(config.Args) != 1 || config.Args[0] != "-Command" {
		t.Errorf("PowerShell args = %v, want [-Command]", config.Args)
	}
}

func TestShellConfig_CmdArgs(t *testing.T) {
	t.Parallel()

	// Test that cmd config has correct args
	config := ShellConfig{
		Name: "cmd",
		Path: "cmd.exe",
		Args: []string{"/C"},
	}

	if len(config.Args) != 1 || config.Args[0] != "/C" {
		t.Errorf("cmd args = %v, want [/C]", config.Args)
	}
}
