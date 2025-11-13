package config

// Config holds the application configuration.
type Config struct {
	Version string `json:"version"` // Hardcoded version, not from TOML
	Foobar  string `json:"foobar"`  // Value from TOML configuration
}

func Create(_ string) error {
	panic("implement me")
}
