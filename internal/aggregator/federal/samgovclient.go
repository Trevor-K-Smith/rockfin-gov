package federal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func CallSamGovAPI() {
	apiKey := "URysJhrDIe4vEkbOZjVqxK6eLCdoDXiU6M9aBg9d"
	postedFrom := "01/01/2025"
	postedTo := "01/08/2025"
	limit := "1"

	baseURL := "https://api.sam.gov/opportunities/v2/search"

	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("postedFrom", postedFrom)
	params.Add("postedTo", postedTo)
	params.Add("limit", limit)

	apiURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP error:", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
