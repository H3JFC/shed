package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"h3jfc/shed/internal/config"
)

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add command that displays the current configuration",
	Long: `The add command prints the current configuration object that was
generated and initialized by the root command. This includes both hardcoded
values and values loaded from the configuration file.`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := config.Get()

		// Print the config object as formatted JSON
		configJSON, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling config: %v\n", err)

			return
		}

		fmt.Println("Current Configuration:")
		fmt.Println(string(configJSON))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
