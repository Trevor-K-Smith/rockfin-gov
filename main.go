package main

import (
	"fmt"
	"os"
	"rockfin-gov/cmd/aggregator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "rockfin-gov",
	Short: "Rockfin Gov CLI",
	Long:  `Rockfin Gov command line interface to fetch and process government opportunities data.`,
}

var aggregatorCmd = aggregator.AggregatorCmd

func initConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	initConfig()
	// fetchCmd.Flags().IntP("limit", "l", 0, "Limit the number of records fetched")
	// samGovCmd.AddCommand(fetchCmd)
	// samGovCmd.AddCommand(samgov.SamgovCmd) // Add PullCmd from cmd/samgov/pull.go
	rootCmd.AddCommand(aggregatorCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
