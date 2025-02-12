package samgov

import (
	"fmt"
	"rockfin-gov/internal/aggregator/federal/samgov"

	db "rockfin-gov/internal/database"

	"github.com/spf13/cobra"
)

// FetchCmd represents the fetch command
var FetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch latest opportunities from sam.gov",
	Long: `This command fetches the latest opportunities from sam.gov
using the samgovclient.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		raw, _ := cmd.Flags().GetBool("raw")

		var data string
		if raw {
			data = samgov.FetchOpportunitiesRaw(limit, "", "")
			fmt.Println(data)
		} else {
			data = samgov.FetchOpportunities(limit, "", "")
			if err := db.ConnectDB(); err != nil {
				fmt.Println("Error connecting to database:", err) // Or handle differently
			} else {
				fmt.Println("Successfully connected to database")
			}
		}
		return nil
	},
}

func init() {
	flags := FetchCmd.Flags()
	flags.IntP("limit", "l", 1000, "Limit the number of records fetched")
	flags.BoolP("raw", "r", false, "Return the raw JSON response")
	SamgovCmd.AddCommand(FetchCmd)
}
