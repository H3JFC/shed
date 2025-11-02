package cmd

import (
	"github.com/spf13/cobra"

	"h3jfc/shed/internal/commands"
	"h3jfc/shed/internal/logger"
)

// initCmd represents the add command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command that displays the current configuration",
	Long:  `init`,
	RunE: func(c *cobra.Command, _ []string) error {
		l := logger.Get()
		l.Info("Initializing shed configuration")
		l.Debug("Debug Init command called")

		if err := commands.Init(c.Context()); err != nil {
			l.Error("Error initializing shed", "error", err)

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
