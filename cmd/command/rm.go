package command

import (
	"errors"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

// RmCmd represents the rm command.
var RmCmd = &cobra.Command{
	Use:   "rm <COMMAND_NAME>",
	Short: "Remove a command from shed",
	Long: `Remove an existing command from shed by name.

This operation is irreversible. The command and all its associated data
will be permanently deleted.

Example:
  # Remove a command
  shed rm list_files

  # Remove a command with verbose output
  shed rm old_command -v`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		commandName := args[0]

		logger.Debug("Removing command", "name", commandName)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		err = s.RemoveCommand(commandName)
		if err != nil {
			if errors.Is(err, store.ErrCommandNotFound) {
				logger.Error("Command not found", "name", commandName)

				return err
			}

			logger.Error("Failed to remove command", "error", err)

			return err
		}

		logger.Info("Command removed successfully", "name", commandName)

		return nil
	},
}
