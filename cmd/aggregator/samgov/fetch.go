package samgov

import (
	"encoding/json"
	"fmt"
	"rockfin-gov/internal/aggregator/federal/samgov"

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

		data := samgov.CallSamGovAPI(limit, "", "")

		if raw {
			fmt.Println(data) // Print raw JSON response
			return nil
		}

		// Parse JSON and format it nicely
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal([]byte(data), &prettyJSON); err != nil {
			return fmt.Errorf("error parsing JSON: %v", err)
		}

		prettyData, err := json.MarshalIndent(prettyJSON, "", "  ")
		if err != nil {
			return fmt.Errorf("error formatting JSON: %v", err)
		}
		fmt.Println(string(prettyData))
		return nil
	},
}

func init() {
	flags := FetchCmd.Flags()
	flags.IntP("limit", "l", 1000, "Limit the number of records fetched")
	flags.BoolP("raw", "r", false, "Return the raw JSON response")
	SamgovCmd.AddCommand(FetchCmd)
}
