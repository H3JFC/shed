package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

// describeCmd represents the describe command.
var describeCmd = &cobra.Command{
	Use:   "describe <COMMAND_NAME>",
	Short: "Display detailed information about a command",
	Long: `Display detailed information about a specific command including its name,
command string, description, parameters, and timestamps.

Example:
  # Describe a command
  shed describe list_files

  # Describe a command with verbose output
  shed describe greet -v`,
	Args: cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		commandName := args[0]

		logger.Debug("Describing command", "name", commandName)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)
			return err
		}

		cmd, err := s.GetCommandByName(commandName)
		if err != nil {
			if errors.Is(err, store.ErrCommandNotFound) {
				logger.Error("Command not found", "name", commandName)
				return err
			}

			logger.Error("Failed to get command", "error", err)
			return err
		}

		var sb strings.Builder

		fmt.Fprintf(&sb, "\nID:          %d\n", cmd.ID)
		fmt.Fprintf(&sb, "Name:        %s\n", cmd.Name)
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
