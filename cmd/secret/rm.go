package secret

import (
	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

// rmCmd represents the remove secret command.
var rmCmd = &cobra.Command{
	Use:   "rm <KEY>",
	Short: "Remove a secret from shed",
	Long: `Remove a secret from shed.

Permanently deletes a secret by its key.

Example:
  # Remove a secret
  shed secret rm github_token

  # Remove multiple secrets (one at a time)
  shed secret rm api_key
  shed secret rm db_password`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		key := args[0]

		logger.Debug("Removing secret", "key", key)

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		err = s.RemoveSecret(key)
		if err != nil {
			logger.Error("Failed to remove secret", "key", key, "error", err)

			return err
		}

		logger.Info("Secret removed successfully", "key", key)

		return nil
	},
}
