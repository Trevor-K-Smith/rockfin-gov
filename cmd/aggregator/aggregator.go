package aggregator

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AggregatorCmd represents the aggregator command
var AggregatorCmd = &cobra.Command{
	Use:   "aggregator",
	Short: "Commands related to data aggregation",
	Long:  `Provides commands to aggregate data from various sources.`,
}

func init() {
	// AggregatorCmd.AddCommand(SamgovCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// AggregatorCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// AggregatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var samGovCmd = &cobra.Command{
	Use:   "samgov",
	Short: "Commands related to SAM.gov API",
	Long:  `Provides commands to interact with the SAM.gov API for government opportunities.`,
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch opportunities from SAM.gov API",
	Long:  `Fetches government opportunities from the SAM.gov API and stores them in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		fmt.Println("Calling samgov API with limit:", limit) // Placeholder - remove actual API call logic
	},
}

func init() {
	AggregatorCmd.AddCommand(samGovCmd)
	fetchCmd.Flags().IntP("limit", "l", 0, "Limit the number of records fetched")
	samGovCmd.AddCommand(fetchCmd)
}
