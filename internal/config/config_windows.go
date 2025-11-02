//go:build windows

package config

var DefaultConfigPaths = []string{
	"%USERPROFILE%",
	"%APPDATA%",
}
