package samgov

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	database "rockfin-gov/internal/database"
	samgovdb "rockfin-gov/internal/database/samgov"
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

			// Unmarshal the response
			var oppResp samgovdb.OpportunitiesResponse
			if err := json.Unmarshal(body, &oppResp); err != nil {
				fmt.Println("Error unmarshaling response:", err)
				return "" // Or perhaps continue and try the next key
			}

			// Create a map to store raw opportunity data by NoticeID
			rawOpportunities := make(map[string]json.RawMessage)
			var rawResp map[string]interface{}
			if err := json.Unmarshal(body, &rawResp); err != nil {
				fmt.Println("Error unmarshaling raw response:", err)
				return ""
			}

			// Check if "opportunitiesData" key exists and is an array
			if oppData, ok := rawResp["opportunitiesData"].([]interface{}); ok {
				for _, opp := range oppData {
					if oppMap, ok := opp.(map[string]interface{}); ok {
						if noticeID, ok := oppMap["noticeId"].(string); ok {
							rawJSON, err := json.Marshal(oppMap)
							if err != nil {
								fmt.Println("Error marshaling raw opportunity:", err)
								return "" // Handle error appropriately
							}
							rawOpportunities[noticeID] = rawJSON
						}
					}
				}
			}

			// Create the samgov database if it doesn't exist
			if err := samgovdb.CreateSamgovDatabase(); err != nil {
				fmt.Println("Error creating samgov database:", err)
				os.Exit(1)
			}

			// Initialize database connection
			database.ConnectDB()

			// Create the table if it doesn't exist
			if err := samgovdb.CreateOpportunitiesTable(database.DB); err != nil {
				fmt.Println("Error creating opportunities table:", err)
				os.Exit(1)
			}

			// Save the opportunities
			// Iterate through oppResp.Opportunities and set the RawData field
			for i := range oppResp.Opportunities {
				if raw, ok := rawOpportunities[oppResp.Opportunities[i].NoticeID]; ok {
					oppResp.Opportunities[i].RawData = raw
				}
			}

			if err := samgovdb.SaveOpportunities(database.DB, oppResp.Opportunities); err != nil {
				fmt.Println("Error saving opportunities:", err)
				os.Exit(1)
			}

			return string(body) // Return the raw body for now, or consider returning a success message
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
