//go:build darwin

package config

var DefaultConfigPaths = []string{
	"$HOME/.shed",
	"$HOME/.config/shed",
}
