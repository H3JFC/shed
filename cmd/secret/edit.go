package secret

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

var editSecretDescription string

const editSecretRequiredArgs = 2

// editCmd represents the edit secret command.
var editCmd = &cobra.Command{
	Use:   "edit <KEY> <VALUE>",
	Short: "Edit an existing secret in shed",
	Long: `Edit an existing secret in shed.

Updates the value and optionally the description of an existing secret.
If --description is not provided, the existing description is preserved.

Example:
  # Update only the value
  shed secret edit github_token ghp_newtoken123

  # Update both value and description
  shed secret edit db_password newpass -d "Updated database password"`,
	Args: cobra.ExactArgs(editSecretRequiredArgs),
	RunE: func(_ *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		logger.Debug("Editing secret", "key", key, "description", editSecretDescription)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		// Get existing secret to preserve description if not provided
		existing, err := s.GetSecretByKey(key)
		if err != nil {
			logger.Error("Secret not found", "key", key, "error", err)

			return fmt.Errorf("failed to update secret: %w", store.ErrSecretNotFound)
		}

		description := existing.Description
		if editSecretDescription != "" {
			description = editSecretDescription
		}

		secret, err := s.UpdateSecret(key, value, description)
		if err != nil {
			if errors.Is(err, store.ErrSecretNotFound) {
				logger.Error("Secret not found", "key", key)

				return err
			}

			if errors.Is(err, store.ErrInvalidCommandName) {
				logger.Error("Invalid secret key", "key", key, "error", err)

				return err
			}

			logger.Error("Failed to update secret", "error", err)

			return err
		}

		logger.Info("Secret updated successfully",
			"id", secret.ID,
			"key", secret.Key,
			"description", secret.Description,
		)

		return nil
	},
}
