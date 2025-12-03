package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

var addDescription string

const addRequiredArgs = 2

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add <COMMAND_NAME> <COMMAND_COMMAND>",
	Short: "Add a new command to shed",
	Long: `Add a new command to shed with a name, description, and command string.

The command string can contain parameters using the {{name|description}} syntax.

Example:
  shed add list_files "ls -la {{path|directory path}}" --description "List files in a directory"
  shed add greet "echo Hello {{name|person's name}}" -d "Greet someone by name"`,
	Args: cobra.ExactArgs(addRequiredArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		commandName := args[0]
		commandCommand := args[1]

		logger.Debug("Adding command", "name", commandName, "command", commandCommand, "description", addDescription)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		cmd, err := s.AddCommand(commandName, commandCommand, addDescription)
		if err != nil {
			if errors.Is(err, store.ErrAlreadyExists) {
				logger.Error("Command already exists", "name", commandName)

				return err
			}

			if errors.Is(err, store.ErrInvalidCommandName) {
				logger.Error("Invalid command name", "name", commandName, "error", err)

				return err
			}

			logger.Error("Failed to add command", "error", err)

			return err
		}

		logger.Info("Command added successfully",
			"id", cmd.ID,
			"name", cmd.Name,
			"command", cmd.Command,
			"description", cmd.Description,
			"parameters", len(cmd.Parameters),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addDescription, "description", "d", "", "Description of the command")
}
