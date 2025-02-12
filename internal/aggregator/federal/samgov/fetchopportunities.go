package samgov

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func FetchOpportunities(limit int, postedFrom, postedTo string) string {
	now := time.Now()

	// Default date range (-7 days to today)
	if postedFrom == "" {
		daysAgo := now.AddDate(0, 0, -7)
		postedFrom = daysAgo.Format("01/02/2006")
	}
	if postedTo == "" {
		postedTo = now.Format("01/02/2006")
	}

	limitStr := strconv.Itoa(limit) // Consistent limit conversion

	baseURL := "https://api.sam.gov/prod/opportunities/v2/search"
	var apiKey string
	var apiURL string

	for attempt := 0; attempt < len(config.SamGov.ApiKeys); attempt++ {
		apiKey = GetCurrentApiKey()

		params := url.Values{}
		params.Add("api_key", apiKey)
		params.Add("postedFrom", postedFrom)
		params.Add("postedTo", postedTo)
		params.Add("limit", limitStr)

		apiURL = baseURL + "?" + params.Encode()

		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error:", err)
			RotateApiKey() // Rotate key on error
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
				os.Exit(1)
			}
			return string(body)
		} else if resp.StatusCode == http.StatusTooManyRequests { // 429
			RotateApiKey()
		} else {
			fmt.Printf("HTTP error: %d, API Key: %s\n", resp.StatusCode, apiKey)
			RotateApiKey() // Rotate key on other errors as well
		}
	}

	fmt.Println("All API keys exhausted.")
	return ""
}
