package cmd

import (
	"errors"

	"github.com/h3jfc/shed/internal/commands"
	"github.com/h3jfc/shed/internal/config"
	"github.com/h3jfc/shed/internal/logger"
	"github.com/spf13/cobra"
)

var ErrShedAlreadyInitialized = errors.New("shed already initialized")

// initCmd represents the add command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command that displays the current configuration",
	Long:  `init`,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		logger.Debug("Starting shed initialization process")

		logger.Debug("Checking for existing shed configuration")
		p, err := config.FindDir()
		if err != nil && !errors.Is(err, config.ErrNoPathFound) {
			return err
		}

		if err == nil && p != "" {
			logger.Info("Shed is already initialized at location", "location", p)
			logger.Debug("Initialization aborted to prevent overwriting existing configuration")

			return ErrShedAlreadyInitialized
		}

		return nil
	},
	RunE: func(c *cobra.Command, _ []string) error {
		logger.Info("Initializing shed configuration")

		if err := commands.Init(c.Context()); err != nil {
			logger.Debug("Error running init command", "error", err)

			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
