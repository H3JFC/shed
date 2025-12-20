//go:build windows

package execute

import (
	"os/exec"
)

// detectShellPlatform detects the appropriate shell for Windows.
func detectShellPlatform() ShellConfig {
	// 1. Check for PowerShell Core (pwsh.exe)
	if path, err := exec.LookPath("pwsh"); err == nil {
		return ShellConfig{
			Name: "pwsh",
			Path: path,
			Args: []string{"-Command"},
		}
	}

	// 2. Check for Windows PowerShell (powershell.exe)
	if path, err := exec.LookPath("powershell"); err == nil {
		return ShellConfig{
			Name: "powershell",
			Path: path,
			Args: []string{"-Command"},
		}
	}

	// 3. Fallback to cmd.exe (always exists on Windows)
	return ShellConfig{
		Name: "cmd",
		Path: "cmd.exe",
		Args: []string{"/C"},
	}
}
