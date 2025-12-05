package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commands in shed",
	Long: `List all stored commands in shed.

Displays the name, command, description, and number of parameters for each command.

Example:
  # List all commands
  shed list

  # List all commands with verbose output
  shed list -v`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		logger.Debug("Listing commands")

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		commands, err := s.ListCommands()
		if err != nil {
			logger.Error("Failed to list commands", "error", err)

			return err
		}

		if len(commands) == 0 {
			logger.Info("No commands found")

			return nil
		}

		logger.Info(fmt.Sprintf("Found %d command(s)", len(commands)))

		for _, cmd := range commands {
			var sb strings.Builder

			fmt.Fprintf(&sb, "\nName:        %s\n", cmd.Name)
			fmt.Fprintf(&sb, "Command:     %s\n", cmd.Command)
			fmt.Fprintf(&sb, "Description: %s\n", cmd.Description)
			fmt.Fprintf(&sb, "Parameters:  %d", len(cmd.Parameters))

			if len(cmd.Parameters) > 0 {
				sb.WriteString("\n  Details:")
				for _, param := range cmd.Parameters {
					if param.Description != "" {
						fmt.Fprintf(&sb, "\n    - %s: %s", param.Name, param.Description)
					} else {
						fmt.Fprintf(&sb, "\n    - %s", param.Name)
					}
				}
			}

			fmt.Fprintf(&sb, "\nCreated:     %s\n", cmd.CreatedAt)
			fmt.Fprintf(&sb, "Updated:     %s", cmd.UpdatedAt)

			logger.Info(sb.String())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
