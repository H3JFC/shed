package cmd

import (
	"log/slog"
	"os"
	"time"

	sf "github.com/samber/slog-formatter"
	"github.com/spf13/cobra"

	"h3jfc/shed/internal/commands"
)

// Usage example - modify your cmd package:
/*
var logMode = ModeFoo // or ModeVerbose, ModeBar

var logger = slog.New(NewCustomHandler(os.Stdout, logMode))

// Or make it configurable via flag:
var logModeFlag string

func init() {
	rootCmd.PersistentFlags().StringVar(&logModeFlag, "log-mode", "foo", "Logging mode: foo, verbose, or bar")
}

// In your initCmd or root command:
func getLogMode(mode string) LogMode {
	switch mode {
	case "verbose":
		return ModeVerbose
	case "bar":
		return ModeBar
	default:
		return ModeFoo
	}
}

var logger = slog.New(NewCustomHandler(os.Stdout, getLogMode(logModeFlag))).
*/
var logger = slog.New(
	sf.NewFormatterHandler(
		sf.UnixTimestampFormatter(time.Millisecond),
	)(slog.NewTextHandler(os.Stdout, nil)),
)

// initCmd represents the add command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command that displays the current configuration",
	Long:  `init`,
	Run: func(c *cobra.Command, _ []string) {
		if err := commands.Init(c.Context()); err != nil {
			logger.Error("Error initializing shed", "error", err)
			os.Exit(1)
		}
		// 1. Detect OS
		// 2. Based on OS check all possible default locatins for existing config and
		//    tell there user a configuration exists and where
		// 3. If non exist, Based on OS offer locations
		// 	a. Windows: %USERPROFILE%\.shed or %APPDATA%\.shed
		// 	b. MacOS: $HOME/.shed or $HOME/.config/shed or /etc/shed
		// 	c. Linux: $HOME/.shed or $HOME/.config/shed or /etc/shed
		// 	d. OR ask the user to provide a custom location
		// 4. Add ShedDirectory to PATH and ask the user to Add to Path based on common shells (bash, zsh, fish, powershell)
		// 5. Create a default config file in the selected location
		// 6. Create a default database file in the selected location
		// 7. Add ShedDir to environment variable SHED_DIR
		// 8. Add ShedDir as a default argument to Cobra commands
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
