//go:build windows

package config

import "os"

var DefaultConfigPaths = []string{
	os.ExpandEnv("%USERPROFILE%"),
	os.ExpandEnv("%APPDATA%"),
}
