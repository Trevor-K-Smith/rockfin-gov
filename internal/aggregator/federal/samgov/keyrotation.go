package samgov

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	currentApiKeyIndex = 0
	config             Config
)

type Config struct {
	SamGov struct {
		ApiKeys []string `yaml:"api_keys"`
	} `yaml:"samgov"`
}

func init() {
	// Load configuration from file
	configPath := "/home/trevorksmith/git/rockfin-gov/config.yaml"
	configFile, err := os.ReadFile(configPath) // Adjust path as necessary
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config:", err)
		os.Exit(1)
	}

	if len(config.SamGov.ApiKeys) == 0 {
		fmt.Println("No API keys found in config")
		os.Exit(1)
	}
}

func RotateApiKey() {
	currentApiKeyIndex = (currentApiKeyIndex + 1) % len(config.SamGov.ApiKeys)
}

func GetCurrentApiKey() string {
	return config.SamGov.ApiKeys[currentApiKeyIndex]
}
