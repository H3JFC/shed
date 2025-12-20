package execute

import "sync"

// ShellConfig represents the configuration for executing commands through a shell.
// It contains the shell name, path, and arguments needed to execute commands.
//
// Example configurations:
//   - bash:       {Name: "bash", Path: "/bin/bash", Args: ["-c"]}
//   - zsh:        {Name: "zsh", Path: "/bin/zsh", Args: ["-c"]}
//   - fish:       {Name: "fish", Path: "/usr/bin/fish", Args: ["-c"]}
//   - PowerShell: {Name: "powershell", Path: "powershell.exe", Args: ["-Command"]}
//   - cmd:        {Name: "cmd", Path: "cmd.exe", Args: ["/C"]}
type ShellConfig struct {
	Name string   // Shell name: "bash", "zsh", "fish", "powershell", "cmd"
	Path string   // Full path to shell executable
	Args []string // Shell arguments for command execution (e.g., ["-c"] for bash)
}

var (
	cachedShell   ShellConfig
	shellMutex    sync.RWMutex
	shellDetected bool
)

// GetShellConfig returns the shell configuration for the current platform.
// The shell is detected once and cached for subsequent calls.
func GetShellConfig() ShellConfig {
	// Fast path: read lock for cached value
	shellMutex.RLock()

	if shellDetected {
		config := cachedShell

		shellMutex.RUnlock()

		return config
	}

	shellMutex.RUnlock()

	// Slow path: detect and cache
	shellMutex.Lock()
	defer shellMutex.Unlock()

	// Double-check after acquiring write lock
	if !shellDetected {
		cachedShell = detectShellPlatform()
		shellDetected = true
	}

	return cachedShell
}

// SetShellConfig allows manual configuration of the shell.
// This is primarily useful for testing.
func SetShellConfig(config ShellConfig) {
	shellMutex.Lock()
	defer shellMutex.Unlock()

	cachedShell = config
	shellDetected = true
}

// ResetShellConfig clears the cached shell configuration,
// forcing re-detection on the next GetShellConfig call.
// This is primarily useful for testing.
func ResetShellConfig() {
	shellMutex.Lock()
	defer shellMutex.Unlock()

	cachedShell = ShellConfig{}
	shellDetected = false
}

// detectShellPlatform is implemented in platform-specific files:
// - shell_darwin.go (macOS)
// - shell_linux.go (Linux)
// - shell_windows.go (Windows)
