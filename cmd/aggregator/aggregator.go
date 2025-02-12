package aggregator

import (
	"rockfin-gov/cmd/aggregator/samgov"

	"github.com/spf13/cobra"
)

// AggregatorCmd represents the aggregator command
var AggregatorCmd = &cobra.Command{
	Use:   "aggregator",
	Short: "Commands related to data aggregation",
	Long:  `Provides commands to aggregate data from various sources.`,
}

func init() {
	AggregatorCmd.AddCommand(samgov.SamgovCmd)
}
