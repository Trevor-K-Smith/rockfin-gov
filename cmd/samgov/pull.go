package samgov

import (
	"fmt"

	"rockfin-gov/internal/aggregator/federal/samgov"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PullCmd represents the pull command
var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest opportunities from sam.gov",
	Long: `This command pulls the latest opportunities from sam.gov
using the samgovclient.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("pull called")
		limit := viper.GetInt("limit")
		data := samgov.CallSamGovAPI(limit)
		fmt.Println(data)
	},
}

func init() {
	SamgovCmd.AddCommand(PullCmd) // Use PullCmd here

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	PullCmd.PersistentFlags().Int("limit", 0, "Limit the number of records to pull") // Use PullCmd here
	viper.BindPFlag("limit", PullCmd.PersistentFlags().Lookup("limit"))              // Use PullCmd here

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
