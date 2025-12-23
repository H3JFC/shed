package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/h3jfc/shed/internal/logger"
	"github.com/h3jfc/shed/internal/store"
	"github.com/spf13/cobra"
)

const (
	cpMinArgs = 2
	cpMaxArgs = 3
)

// CpCmd represents the cp command.
var CpCmd = &cobra.Command{
	Use:   "cp <COMMAND_SRC_NAME> <COMMAND_DEST_NAME> [jsonValueParams]",
	Short: "Copy a command with optional parameter value substitutions",
	Long: `Copy an existing command to a new name, optionally providing values for parameters.

Parameters should be provided as a JSON object in the form {"param":"value"}.
The command will substitute these values in the copied command and remove those
parameters from the new command.

Example:
  # Copy without parameter substitution
  shed cp list_files list_home_files

  # Copy with single parameter substitution
  shed cp list_files list_home_files '{"path":"/home/user"}'

  # Copy with multiple parameter substitutions
  shed cp greet greet_john '{"name":"John","title":"Mr."}'`,
	Args: cobra.RangeArgs(cpMinArgs, cpMaxArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		srcName := args[0]
		destName := args[1]

		jsonValueParams := "{}"
		if len(args) == cpMaxArgs {
			jsonValueParams = args[2]
		}

		logger.Debug("Copying command",
			"src", srcName,
			"dest", destName,
			"params", jsonValueParams,
		)

		// Validate JSON format
		if err := validateJSON(jsonValueParams); err != nil {
			logger.Error("Invalid JSON parameter format", "error", err)

			return fmt.Errorf("invalid JSON format: %w", err)
		}

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		cmd, err := s.CopyCommand(srcName, destName, jsonValueParams)
		if err != nil {
			if errors.Is(err, store.ErrCommandNotFound) {
				logger.Error("Source command not found", "name", srcName)

				return err
			}

			if errors.Is(err, store.ErrAlreadyExists) {
				logger.Error("Destination command already exists", "name", destName)

				return err
			}

			if errors.Is(err, store.ErrInvalidCommandName) {
				logger.Error("Invalid destination command name", "name", destName, "error", err)

				return err
			}

			if errors.Is(err, store.ErrParsingValueParams) {
				logger.Error("Failed to parse parameter values", "error", err)

				return err
			}

			logger.Error("Failed to copy command", "error", err)

			return err
		}

		logger.Info("Command copied successfully",
			"id", cmd.ID,
			"name", cmd.Name,
			"command", cmd.Command,
			"description", cmd.Description,
			"parameters", len(cmd.Parameters),
		)

		return nil
	},
}

// validateJSON checks if a string is valid JSON object format.
func validateJSON(jsonStr string) error {
	var m map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}
