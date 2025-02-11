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

func CallSamGovAPI(limit int) string {
	now := time.Now()
	dayBeforeYesterday := now.AddDate(0, 0, -25)
	yesterday := now.AddDate(0, 0, -20)
	postedFrom := dayBeforeYesterday.Format("01/02/2006")
	postedTo := yesterday.Format("01/02/2006")

	limitStr := "1000" // Default limit
	if limit != 0 {
		limitStr = strconv.Itoa(limit)
	}

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

func CallSamGovAPIWithLimit(limit int) string {
	now := time.Now()
	daysAgo := now.AddDate(0, 0, -7)
	postedFrom := daysAgo.Format("01/02/2006")
	postedTo := now.Format("01/02/2006")
	limitStr := fmt.Sprintf("%d", limit)

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
	os.Exit(1) // Exit if all keys are tried
	return ""
}
