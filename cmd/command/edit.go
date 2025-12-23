package command

import (
	"errors"

	"github.com/h3jfc/shed/internal/logger"
	"github.com/h3jfc/shed/internal/store"
	"github.com/spf13/cobra"
)

var (
	editDescription string
	editName        string
)

const (
	editMinArgs = 2
	editMaxArgs = 3
)

// EditCmd represents the edit command.
var EditCmd = &cobra.Command{
	Use:   "edit <COMMAND_NAME> [flags] <CLI_COMMAND> [jsonValueParams]",
	Short: "Edit an existing command in shed",
	Long: `Edit an existing command in shed by updating its name, description, command string, or parameters.

The command string can contain parameters using the {{name|description}} syntax.
You can optionally provide JSON value parameters to hydrate/substitute specific parameters.

Examples:
  # Edit command string only
  shed edit list_files "ls -lah {{path|directory path}}"

  # Edit with new description
  shed edit list_files --description "List all files" "ls -la {{path}}"

  # Edit and rename command
  shed edit list_files --name show_files "ls -la {{path}}"

  # Edit and hydrate a parameter
  shed edit api_call "curl -XGET {{url}} -H {{auth}}" '{"url":"https://api.example.com"}'

  # Edit everything at once
  shed edit old_name --name new_name --description "New description" "new command {{param}}" '{"other":"value"}'`,
	Args: cobra.RangeArgs(editMinArgs, editMaxArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		commandName := args[0]
		commandCommand := args[1]

		jsonValueParams := ""
		if len(args) == editMaxArgs {
			jsonValueParams = args[2]
		}

		logger.Debug("Editing command",
			"name", commandName,
			"command", commandCommand,
			"newName", editName,
			"description", editDescription,
			"jsonValueParams", jsonValueParams,
		)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		// Get the existing command
		existingCmd, err := s.GetCommandByName(commandName)
		if err != nil {
			if errors.Is(err, store.ErrCommandNotFound) {
				logger.Error("Command not found", "name", commandName)

				return err
			}

			logger.Error("Failed to get command", "error", err)

			return err
		}

		// Determine the new name (use existing if not provided)
		newName := commandName
		if editName != "" {
			newName = editName
		}

		// Determine the new description (use existing if not provided)
		newDescription := existingCmd.Description
		if editDescription != "" {
			newDescription = editDescription
		}

		// Update the command
		updatedCmd, err := s.UpdateCommand(
			existingCmd.ID,
			newName,
			commandCommand,
			newDescription,
			existingCmd.Parameters,
			jsonValueParams,
		)
		if err != nil {
			if errors.Is(err, store.ErrInvalidCommandName) {
				logger.Error("Invalid command name", "name", newName, "error", err)

				return err
			}

			if errors.Is(err, store.ErrParsingValueParams) {
				logger.Error("Invalid JSON value parameters", "error", err)

				return err
			}

			logger.Error("Failed to update command", "error", err)

			return err
		}

		logger.Info("Command updated successfully",
			"id", updatedCmd.ID,
			"name", updatedCmd.Name,
			"command", updatedCmd.Command,
			"description", updatedCmd.Description,
			"parameters", len(updatedCmd.Parameters),
		)

		return nil
	},
}

func init() {
	EditCmd.Flags().StringVarP(&editDescription, "description", "d", "", "New description for the command")
	EditCmd.Flags().StringVarP(&editName, "name", "n", "", "New name for the command")
}
