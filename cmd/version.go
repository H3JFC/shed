package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/logger"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	Long: `Display the current version and build information for shed.

Shows the version number and commit hash used to build this binary.

Example:
  shed version`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		logger.Info(fmt.Sprintf("shed version %s (commit: %s)\n", Version, Commit))
	},
}
