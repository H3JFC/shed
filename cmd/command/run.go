package command

import (
	"encoding/json"
	"fmt"

	"github.com/h3jfc/shed/internal/execute"
	"github.com/h3jfc/shed/internal/logger"
	"github.com/h3jfc/shed/internal/store"
	"github.com/h3jfc/shed/lib/brackets"
	"github.com/spf13/cobra"
)

const (
	maxRunArgs = 2
)

// RunCmd represents the run command.
var RunCmd = &cobra.Command{
	Use:   "run <COMMAND_NAME> [jsonValueParams]",
	Short: "Run a stored command by name",
	Long: `Run a stored command by name, with optional parameter values.

Parameters should be provided as a JSON object in the form {"param":"value"}.
If a command requires parameters and they are not provided, an error will be returned.

Secrets (parameters starting with !) are automatically fetched from the secrets store
and substituted into the command before execution.

Examples:
  # Run a command without parameters
  shed run list_files

  # Run a command with a single parameter
  shed run list_files '{"path":"/home/user"}'

  # Run a command with multiple parameters
  shed run deploy '{"environment":"production","version":"1.2.3"}'

  # List available commands
  shed list`,
	Args: cobra.RangeArgs(1, maxRunArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		commandName := args[0]

		jsonValueParams := "{}"
		if len(args) == maxRunArgs {
			jsonValueParams = args[1]
		}

		logger.Debug("Running command", "name", commandName, "params", jsonValueParams)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		// Get the command
		cmd, err := s.GetCommandByName(commandName)
		if err != nil {
			logger.Error("Failed to get command", "name", commandName, "error", err)

			return err
		}

		logger.Debug("Retrieved command",
			"name", cmd.Name,
			"command", cmd.Command,
			"parameters", len(cmd.Parameters),
		)

		// Validate JSON format
		if err := validateJSON(jsonValueParams); err != nil {
			logger.Error("Invalid JSON parameter format", "error", err)

			return fmt.Errorf("invalid JSON format: %w", err)
		}

		// Parse the command to extract secrets
		parsed, err := brackets.Parse(cmd.Command)
		if err != nil {
			logger.Error("Failed to parse command", "error", err)

			return fmt.Errorf("failed to parse command: %w", err)
		}

		// Parse the provided parameters
		var paramMap map[string]string
		if err := json.Unmarshal([]byte(jsonValueParams), &paramMap); err != nil {
			logger.Error("Failed to parse parameters", "error", err)

			return fmt.Errorf("failed to parse parameters: %w", err)
		}

		// Fetch secrets and add them to parameter map
		for _, secret := range *parsed.Secrets {
			secretValue, err := s.GetSecretByKey(secret.Key)
			if err != nil {
				logger.Error("Failed to get secret", "key", secret.Key, "error", err)

				return fmt.Errorf("failed to get secret %s: %w", secret.Key, err)
			}
			// Secret parameters are prefixed with ! in the command string
			paramMap["!"+secret.Key] = secretValue.Value
			logger.Debug("Loaded secret", "key", secret.Key)
		}

		// Convert parameter map back to JSON for hydration
		updatedParams, err := json.Marshal(paramMap)
		if err != nil {
			logger.Error("Failed to marshal parameters", "error", err)

			return fmt.Errorf("failed to marshal parameters: %w", err)
		}

		// Hydrate the command with parameter values
		hydratedCmd, err := brackets.HydrateStringFromJSON(cmd.Command, string(updatedParams))
		if err != nil {
			logger.Error("Failed to hydrate command", "error", err)

			return fmt.Errorf("failed to hydrate command: %w", err)
		}

		logger.Debug("Hydrated command", "command", hydratedCmd)
		logger.Info("Executing command", "name", cmd.Name)

		// Execute the command
		if err := execute.Run(hydratedCmd); err != nil {
			logger.Error("Command execution failed", "error", err)

			return fmt.Errorf("command execution failed: %w", err)
		}

		logger.Info("Command executed successfully", "name", cmd.Name)

		return nil
	},
}
