package secret

import (
	"github.com/spf13/cobra"
)

// Cmd represents the parent secret command.
var Cmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets in shed",
	Long: `Manage secrets in shed.

Secrets are used to store sensitive information like passwords, API keys, and tokens.
They can be referenced in commands using the {{!key}} syntax.

Available commands:
  add     Add a new secret
  list    List all secrets
  edit    Edit an existing secret
  rm      Remove a secret`,
}

// Init registers all secret subcommands with the parent command.
func Init() *cobra.Command {
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(editCmd)
	Cmd.AddCommand(rmCmd)

	addCmd.Flags().StringVarP(&addSecretDescription, "description", "d", "", "Description of the secret")
	editCmd.Flags().StringVarP(&editSecretDescription, "description", "d", "", "New description for the secret")

	return Cmd
}
