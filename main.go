package main

import (
	"fmt"
	"os"
	"rockfin-gov/cmd/samgov"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "rockfin-gov",
	Short: "Rockfin Gov CLI",
	Long:  `Rockfin Gov command line interface to fetch and process government opportunities data.`,
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
	fetchCmd.Flags().IntP("limit", "l", 0, "Limit the number of records fetched")
	samGovCmd.AddCommand(fetchCmd)
	samGovCmd.AddCommand(samgov.PullCmd) // Add PullCmd from cmd/samgov/pull.go
	rootCmd.AddCommand(samGovCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
