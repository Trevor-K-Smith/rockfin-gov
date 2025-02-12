package samgov

import (
	"github.com/spf13/cobra"
)

// SamgovCmd represents the samgov command
var SamgovCmd = &cobra.Command{
	Use:   "aggregator samgov",
	Short: "Commands related to sam.gov",
	Long:  `This command group contains subcommands for interacting with sam.gov`,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// SamgovCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// SamgovCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
