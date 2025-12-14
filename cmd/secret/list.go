package secret

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
	"h3jfc/shed/internal/store"
)

// listCmd represents the list secrets command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in shed",
	Long: `List all stored secrets in shed.

Displays the key, description, and timestamps for each secret.
Note: Secret values are not displayed for security reasons.

Example:
  # List all secrets
  shed secret list

  # List all secrets with verbose output
  shed secret list -v`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		logger.Debug("Listing secrets")

		s, err := store.NewStoreFromConfig()
		if err != nil {
			logger.Error("Failed to initialize store", "error", err)

			return err
		}

		secrets, err := s.ListSecrets()
		if err != nil {
			logger.Error("Failed to list secrets", "error", err)

			return err
		}

		if len(secrets) == 0 {
			logger.Info("No secrets found")

			return nil
		}

		logger.Info(fmt.Sprintf("Found %d secret(s)", len(secrets)))

		for _, secret := range secrets {
			var sb strings.Builder

			fmt.Fprintf(&sb, "\nKey:         %s\n", secret.Key)
			fmt.Fprintf(&sb, "Description: %s\n", secret.Description)
			fmt.Fprintf(&sb, "Created:     %s\n", secret.CreatedAt)
			fmt.Fprintf(&sb, "Updated:     %s", secret.UpdatedAt)

			logger.Info(sb.String())
		}

		return nil
	},
}
