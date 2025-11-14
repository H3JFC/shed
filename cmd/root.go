package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"h3jfc/shed/internal/logger"
)

var Commit = "NOT SET"

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "shed",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRunE: func(c *cobra.Command, _ []string) error {
		ll := "message-level"
		if c.Flags().Lookup("verbose").Value.String() == "true" {
			ll = "verbose"
		}

		logger.New(logger.ModeFromString(ll))

		if !isInitCommand(c) {
			initConfig()

			if err := viper.BindPFlags(c.Flags()); err != nil {
				logger.Debug("Error binding flags to viper config", "error", err)
				logger.Error("Error with flag configuration")

				return err
			}
		}

		return nil
	},

	Run: func(_ *cobra.Command, _ []string) {
		// Do Stuff Here
		logger.Warn("Shed is in early development. Use at your own risk!", "commit", viper.GetString("commit"))
		fmt.Println("Shed - A tool for managing your projects")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing command", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("shed-dir", os.Getenv("SHED_DIR"), "Path to the Shed configuration directory")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	// Initialize the configuration system
	logger.Info("initconfig")

	if err := Init(); err != nil {
		logger.Error("Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

// Init initializes the configuration system.
func Init() error {
	// Set up Viper
	viper.SetConfigName("shed")
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("SHED")
	viper.AutomaticEnv()

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create .shed directory if it doesn't exist.
	shedDir := filepath.Join(homeDir, ".shed")
	viper.AddConfigPath(shedDir)

	// Set defaults
	viper.SetDefault("commit", Commit)

	// Read config file (it's okay if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

func isInitCommand(cmd *cobra.Command) bool {
	if cmd.CalledAs() == initCmd.Name() {
		return true
	}

	if initCmd.Aliases != nil {
		return slices.Contains(initCmd.Aliases, cmd.CalledAs())
	}

	return false
}
