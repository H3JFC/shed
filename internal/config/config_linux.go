//go:build linux

package config

import "os"

var DefaultConfigPaths = []string{
	os.ExpandEnv("$HOME/.shed"),
	os.ExpandEnv("$HOME/.config/shed"),
}
