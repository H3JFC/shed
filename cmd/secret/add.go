package secret

import (
	"errors"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

var addSecretDescription string

const addSecretRequiredArgs = 2

// addCmd represents the add secret command.
var addCmd = &cobra.Command{
	Use:   "add <KEY> <VALUE>",
	Short: "Add a new secret to shed",
	Long: `Add a new secret to shed with a key, value, and optional description.

Secrets are used to store sensitive information like passwords, API keys, and tokens.
They can be referenced in commands using the {{!key}} syntax.

Example:
  shed secret add github_token ghp_abc123xyz --description "GitHub API token"
  shed secret add db_password mysecretpass -d "Database password"`,
	Args: cobra.ExactArgs(addSecretRequiredArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		logger.Debug("Adding secret", "key", key, "description", addSecretDescription)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		secret, err := s.AddSecret(key, value, addSecretDescription)
		if err != nil {
			if errors.Is(err, store.ErrAlreadyExists) {
				logger.Error("Secret already exists", "key", key)

				return err
			}

			if errors.Is(err, store.ErrInvalidCommandName) {
				logger.Error("Invalid secret key", "key", key, "error", err)

				return err
			}

			logger.Error("Failed to add secret", "error", err)

			return err
		}

		logger.Info("Secret added successfully",
			"id", secret.ID,
			"key", secret.Key,
			"description", secret.Description,
		)

		return nil
	},
}
