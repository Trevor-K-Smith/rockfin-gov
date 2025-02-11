package federal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SamGov struct {
		ApiKeys []string `yaml:"api_keys"`
	} `yaml:"samgov"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

var (
	currentApiKeyIndex = 0
	config             Config
)

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

func CallSamGovAPI(limit int) string {
	now := time.Now()
	daysAgo := now.AddDate(0, 0, -7)
	postedFrom := daysAgo.Format("01/02/2006")
	postedTo := now.Format("01/02/2006")

	limitStr := "1000" // Default limit
	if limit != 0 {
		limitStr = strconv.Itoa(limit)
	}

	baseURL := "https://api.sam.gov/opportunities/v2/search"
	var apiKey string
	var apiURL string

	for attempt := 0; attempt < len(config.SamGov.ApiKeys); attempt++ {
		apiKey = config.SamGov.ApiKeys[currentApiKeyIndex]

		params := url.Values{}
		params.Add("api_key", apiKey)
		params.Add("postedFrom", postedFrom)
		params.Add("postedTo", postedTo)
		params.Add("limit", limitStr)

		apiURL = baseURL + "?" + params.Encode()

		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error:", err)
			// Rotate key and retry immediately on network errors
			rotateApiKey()
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
				os.Exit(1)
			}
			db, err := connectDB()
			if err != nil {
				fmt.Println("Error connecting to DB:", err)
				os.Exit(1)
			}
			defer db.Close()

			err = insertOpportunityData(db, body)
			if err != nil {
				fmt.Println("Error inserting data:", err)
				os.Exit(1)
			}
			return string(body)
		} else if resp.StatusCode == http.StatusTooManyRequests { //429
			rotateApiKey()
			// Add a delay before retrying
		} else {
			fmt.Printf("HTTP error: %d, API Key: %s\n", resp.StatusCode, apiKey)
			// Rotate on other errors as well, in case it's an issue with a specific key
			rotateApiKey()
		}
	}
	fmt.Println("All API keys exhausted.")
	os.Exit(1) // Exit if all keys have been tried
	return ""
}

func CallSamGovAPIWithLimit(limit int) string {
	now := time.Now()
	daysAgo := now.AddDate(0, 0, -7)
	postedFrom := daysAgo.Format("01/02/2006")
	postedTo := now.Format("01/02/2006")
	limitStr := fmt.Sprintf("%d", limit)

	baseURL := "https://api.sam.gov/opportunities/v2/search"
	var apiKey string
	var apiURL string

	for attempt := 0; attempt < len(config.SamGov.ApiKeys); attempt++ {
		apiKey = config.SamGov.ApiKeys[currentApiKeyIndex]

		params := url.Values{}
		params.Add("api_key", apiKey)
		params.Add("postedFrom", postedFrom)
		params.Add("postedTo", postedTo)
		params.Add("limit", limitStr)

		apiURL = baseURL + "?" + params.Encode()

		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error:", err)
			// Rotate key and retry immediately on network errors
			rotateApiKey()
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
		} else if resp.StatusCode == http.StatusTooManyRequests { //429
			rotateApiKey()
			// Add a delay before retrying
		} else {
			fmt.Printf("HTTP error: %d, API Key: %s\n", resp.StatusCode, apiKey)
			// Rotate on other errors as well, in case it's an issue with a specific key
			rotateApiKey()
		}
	}
	fmt.Println("All API keys exhausted.")
	os.Exit(1) // Exit if all keys have been tried
	return ""
}

func rotateApiKey() {
	currentApiKeyIndex = (currentApiKeyIndex + 1) % len(config.SamGov.ApiKeys)
}

func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertOpportunityData(db *sql.DB, body []byte) error {
	// 1. Parse the JSON response
	var data map[string]interface{} // Use a generic map for initial parsing
	err := json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS opportunities (
	        id SERIAL PRIMARY KEY,
	        total_records INTEGER NULL,
	        limit_value INTEGER NULL,
	        offset_value INTEGER NULL,
	        title TEXT NULL,
	        solicitation_number TEXT NULL,
	        full_parent_path_name TEXT NULL,
	        full_parent_path_code TEXT NULL,
	        posted_date TIMESTAMP NULL,
	        type TEXT NULL,
	        base_type TEXT NULL,
	        archive_type TEXT NULL,
	        archive_date TEXT NULL,
	        set_aside TEXT NULL,
	        set_aside_code TEXT NULL,
	        reponse_dead_line TEXT NULL,
	        naics_code TEXT NULL,
	        classification_code TEXT NULL,
	        active TEXT NULL,
	        award_number TEXT NULL,
	        award_amount NUMERIC NULL,
	        award_date TIMESTAMP NULL,
	        awardee_name TEXT NULL,
	        awardee_uei_sam TEXT NULL,
	        awardee_street_address TEXT NULL,
	        awardee_street_address2 TEXT NULL,
	        awardee_city_code TEXT NULL,
	        awardee_city_name TEXT NULL,
	        awardee_state_code TEXT NULL,
	        awardee_state_name TEXT NULL,
	        awardee_country_code TEXT NULL,
	        awardee_country_name TEXT NULL,
	        awardee_zip TEXT NULL,
	        point_of_contact JSONB NULL,
	        description TEXT NULL,
	        additional_info JSONB NULL,
	        organization_type TEXT NULL,
	        office_address_city TEXT NULL,
	        office_address_state TEXT NULL,
	        office_address_zip TEXT NULL,
	        pop_street_address TEXT NULL,
	        pop_street_address2 TEXT NULL,
	        pop_city_code TEXT NULL,
	        pop_city_name TEXT NULL,
	        pop_state_code TEXT NULL,
	        pop_state_name TEXT NULL,
	        pop_country_code TEXT NULL,
	        pop_country_name TEXT NULL,
	        pop_zip TEXT NULL,
	        additional_info_link TEXT NULL,
	        ui_link TEXT NULL,
	        links JSONB NULL,
	        resource_links TEXT[] NULL
	    )
	`)
	if err != nil {
		return fmt.Errorf("failed to create table opportunities")
	}

	// 2. Extract relevant data
	opportunitiesData, ok := data["opportunitiesData"].([]interface{})
	if !ok {
		return fmt.Errorf("opportunitiesData not found or not an array")
	}

	foundCount := 0
	addedCount := 0

	for _, opportunity := range opportunitiesData {
		opportunityMap, ok := opportunity.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid data: each opportunity must be a map")
		}

		// Helper function to extract string values safely
		getString := func(key string) string {
			if val, ok := opportunityMap[key]; ok {
				if strVal, ok := val.(string); ok {
					return strVal
				} else {
					return fmt.Sprintf("%v", val)
				}
			}
			return "" // Handle missing field gracefully
		}

		// Helper function to extract integer values safely
		getInt := func(key string) int {
			if val, ok := opportunityMap[key]; ok {
				if intVal, ok := val.(float64); ok { // JSON numbers are float64
					return int(intVal)
				} else if val == nil {
					return 0 // Handle nil values gracefully
				} else {
					fmt.Printf("Warning: Invalid data type for %s, expected integer, got %T\n", key, val)
				}
			}
			return 0 // Handle missing field gracefully
		}

		// Helper function to extract numeric values safely
		getNumeric := func(key string) float64 {
			if val, ok := opportunityMap[key]; ok {
				if numVal, ok := val.(float64); ok {
					return numVal
				} else {
					fmt.Printf("Warning: Invalid data type for %s, expected numeric\n", key)
				}
			}
			return 0 // Handle missing field gracefully
		}

		// Helper function to extract time values safely
		getTime := func(key string) time.Time {
			if val, ok := opportunityMap[key]; ok {
				if timeStr, ok := val.(string); ok {
					// Try parsing different formats
					timeValue, err := time.Parse("2006-01-02", timeStr)
					if err != nil {
						timeValue, err = time.Parse("2006-01-02 15:04:05", timeStr)
						if err != nil {
							fmt.Printf("Warning: Could not parse time for %s: %v\n", key, err)
							return time.Time{} // Return zero time on error
						}
					}
					return timeValue
				} else {
					fmt.Printf("Warning: Invalid data type for %s, expected string\n", key)
				}
			}
			return time.Time{} // Return zero time on error
		}

		// Extract top-level fields
		totalRecords := getInt("totalRecords")
		limitValue := getInt("limit")
		offsetValue := getInt("offset")
		title := getString("title")
		solicitationNumber := getString("solicitationNumber")
		fullParentPathName := getString("fullParentPathName")
		fullParentPathCode := getString("fullParentPathCode")
		postedDate := getTime("postedDate")
		opportunityType := getString("type")
		baseType := getString("baseType")
		archiveType := getString("archiveType")
		archiveDate := getTime("archiveDate") // keep this line
		setAside := getString("setAside")
		setAsideCode := getString("setAsideCode")
		reponseDeadLine := getString("reponseDeadLine")
		naicsCode := getString("naicsCode")
		classificationCode := getString("classificationCode")
		active := getString("active")
		description := getString("description")
		additionalInfoLink := getString("additionalInfoLink")
		uiLink := getString("uiLink")
		organizationType := getString("organizationType")

		// Extract award data
		awardData, _ := opportunityMap["data"].(map[string]interface{})
		var awardNumber string
		var awardAmount float64
		var awardDate time.Time
		var awardeeName string
		var awardeeUeiSAM string
		var awardeeStreetAddress string
		var awardeeStreetAddress2 string
		var awardeeCityCode string
		var awardeeCityName string
		var awardeeStateCode string
		var awardeeStateName string
		var awardeeCountryCode string
		var awardeeCountryName string
		var awardeeZip string

		if awardData != nil {
			award, _ := awardData["award"].(map[string]interface{})
			if award != nil {
				awardNumber = getString("number")
				awardAmount = getNumeric("amount")
				awardDate = getTime("date")

				awardee, _ := award["awardee"].(map[string]interface{})
				if awardee != nil {
					awardeeName = getString("name")
					awardeeUeiSAM = getString("ueiSAM")

					location, _ := awardee["location"].(map[string]interface{})
					if location != nil {
						awardeeStreetAddress = getString("streetAddress")
						awardeeStreetAddress2 = getString("streetAddress2")

						city, _ := location["city"].(map[string]interface{})
						if city != nil {
							awardeeCityCode = getString("code")
							awardeeCityName = getString("name")
						}

						state, _ := location["state"].(map[string]interface{})
						if state != nil {
							awardeeStateCode = getString("code")
							awardeeStateName = getString("name")
						}

						country, _ := location["country"].(map[string]interface{})
						if country != nil {
							awardeeCountryCode = getString("code")
							awardeeCountryName = getString("name")
						}

						awardeeZip = getString("zip")
					}
				}
			}
		}

		// Extract office address data
		officeAddressData, _ := opportunityMap["data"].(map[string]interface{})
		var officeAddressCity string
		var officeAddressState string
		var officeAddressZip string

		if officeAddressData != nil {
			officeAddress, _ := officeAddressData["officeAddress"].(map[string]interface{})
			if officeAddress != nil {
				officeAddressCity = getString("city")
				officeAddressState = getString("state")
				officeAddressZip = getString("zip")
			}
		}

		// Extract place of performance data
		placeOfPerformanceData, _ := opportunityMap["data"].(map[string]interface{})
		var popStreetAddress string
		var popStreetAddress2 string
		var popCityCode string
		var popCityName string
		var popStateCode string
		var popStateName string
		var popCountryCode string
		var popCountryName string
		var popZip string

		if placeOfPerformanceData != nil {
			placeOfPerformance, _ := placeOfPerformanceData["placeOfPerformance"].(map[string]interface{})
			if placeOfPerformance != nil {
				popStreetAddress = getString("streetAddress")
				popStreetAddress2 = getString("streetAddress2")

				city, _ := placeOfPerformance["city"].(map[string]interface{})
				if city != nil {
					popCityCode = getString("code")
					popCityName = getString("name")
				}

				state, _ := placeOfPerformance["state"].(map[string]interface{})
				if state != nil {
					popStateCode = getString("code")
					popStateName = getString("name")
				}

				country, _ := placeOfPerformance["country"].(map[string]interface{})
				if country != nil {
					popCountryCode = getString("code")
					popCountryName = getString("name")
				}

				popZip = getString("zip")
			}
		}

		// Extract point of contact data
		pointOfContactData, ok := opportunityMap["data"].(map[string]interface{})
		var pointOfContact interface{}
		if ok {
			pointOfContactDataList, ok := pointOfContactData["pointOfContact"].([]interface{})
			if ok && len(pointOfContactDataList) > 0 {
				pointOfContact = pointOfContactDataList
			} else {
				pointOfContact = nil
			}
		} else {
			pointOfContact = nil
		}

		// Extract additional info data
		additionalInfoData, ok := opportunityMap["data"].(map[string]interface{})
		var additionalInfo interface{}
		if ok {
			additionalInfoDataList, ok := additionalInfoData["additionalInfo"].([]interface{})
			if ok && len(additionalInfoDataList) > 0 {
				additionalInfo = additionalInfoDataList
			} else {
				additionalInfo = nil
			}
		} else {
			additionalInfo = nil
		}

		// Extract links data
		linksData, _ := opportunityMap["links"].([]interface{})
		var linksJSON []byte
		if len(linksData) > 0 {
			linksJSON, err = json.Marshal(linksData)
			if err != nil {
				fmt.Printf("Warning: Could not marshal links data: %v\n", err)
				linksJSON = []byte("null")
			}
		}

		// Extract resource links data
		resourceLinksData, _ := opportunityMap["resourceLinks"].([]interface{})
		var resourceLinks []string
		if len(resourceLinksData) > 0 {
			resourceLinks = make([]string, len(resourceLinksData))
			for i, link := range resourceLinksData {
				resourceLinks[i] = link.(string)
			}
		}

		// Convert resourceLinks to a PostgreSQL array literal
		resourceLinksString := "{}"
		if len(resourceLinks) > 0 {
			quotedLinks := make([]string, len(resourceLinks))
			for i, link := range resourceLinks {
				quotedLinks[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(link, "\"", "\\\""))
			}
			resourceLinksString = "{" + strings.Join(quotedLinks, ",") + "}"
		}

		// Check if opportunity already exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM opportunities WHERE solicitation_number = $1)", solicitationNumber).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking for existing opportunity: %w", err)
		}

		if exists {
			foundCount++
			//fmt.Printf("Opportunity with solicitation number %s already exists, skipping insertion\n", solicitationNumber)
			continue // Skip to the next opportunity
		}

		// 3. Prepare the SQL INSERT statement
		stmt, err := db.Prepare(`
			INSERT INTO opportunities (
				total_records, limit_value, offset_value, title, solicitation_number,
				full_parent_path_name, full_parent_path_code, posted_date, type, base_type,
				archive_type, archive_date, set_aside, set_aside_code, reponse_dead_line,
				naics_code, classification_code, active, award_number, award_amount,
				award_date, awardee_name, awardee_uei_sam, awardee_street_address, awardee_street_address2,
				awardee_city_code, awardeeCityName, awardee_state_code, awardee_state_name, awardee_country_code,
				awardee_country_name, awardee_zip, point_of_contact, description, additional_info,
				organization_type, office_address_city, office_address_state, office_address_zip, pop_street_address,
				pop_street_address2, pop_city_code, pop_city_name, pop_state_code, pop_state_name,
				pop_country_code, pop_country_name, pop_zip, additional_info_link, ui_link,
				links, resource_links
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20,
				$21, $22, $23, $24, $25,
				$26, $27, $28, $29, $30,
				$31, $32, $33, $34, $35,
				$36, $37, $38, $39, $40,
				$41, $42, $43, $44, $45,
				$46, $47, $48, $49, $50,
				$51, $52
			)
		`)
		if err != nil {
			return fmt.Errorf("error preparing statement: %w", err)
		}
		defer stmt.Close()

		// 4. Execute the INSERT statement
		_, err = stmt.Exec(
			totalRecords, limitValue, offsetValue, title, solicitationNumber,
			fullParentPathName, fullParentPathCode, postedDate, opportunityType, baseType,
			archiveType, archiveDate, setAside, setAsideCode, reponseDeadLine,
			naicsCode, classificationCode, active, awardNumber, awardAmount,
			awardDate, awardeeName, awardeeUeiSAM, awardeeStreetAddress, awardeeStreetAddress2,
			awardeeCityCode, awardeeCityName, awardeeStateCode, awardeeStateName, awardeeCountryCode,
			awardeeCountryName, awardeeZip, pointOfContact, description, additionalInfo,
			organizationType, officeAddressCity, officeAddressState, officeAddressZip, popStreetAddress,
			popStreetAddress2, popCityCode, popCityName, popStateCode, popStateName,
			popCountryCode, popCountryName, popZip, additionalInfoLink, uiLink,
			linksJSON, resourceLinksString,
		)
		if err != nil {
			return fmt.Errorf("error executing statement: %w", err)
		}
		addedCount++
	}

	return nil
}
