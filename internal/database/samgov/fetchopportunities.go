package samgov

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Opportunity struct {
	NoticeID           string          `json:"noticeId" gorm:"primaryKey"`
	SolicitationNumber string          `json:"solicitationNumber"`
	RawData            json.RawMessage `json:"rawData"`
}

type OpportunitiesResponse struct {
	TotalRecords  int           `json:"totalRecords"`
	Limit         int           `json:"limit"`
	Offset        int           `json:"offset"`
	Opportunities []Opportunity `json:"opportunitiesData"`
}

func CreateSamgovDatabase() error {
	var config DBConfig
	err := viper.UnmarshalKey("database", &config)
	if err != nil {
		return fmt.Errorf("failed to decode database config: %v", err)
	}

	// First, connect to the default 'postgres' database to check if 'samgov' exists
	postgresDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, "postgres") // Connect to "postgres" db
	postgresDB, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	var count int64
	result := postgresDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", config.Name).Scan(&count)
	if result.Error != nil {
		return fmt.Errorf("failed to check if database exists: %w", result.Error)
	}

	if count == 0 {
		// Create the 'samgov' database if it doesn't exist
		result := postgresDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Name))
		if result.Error != nil {
			return fmt.Errorf("failed to create database %s: %v", config.Name, result.Error)
		}
	}
	return nil
}

func CreateOpportunitiesTable(conn *gorm.DB) error {
	err := conn.AutoMigrate(&Opportunity{})
	return err
}

func SaveOpportunity(conn *gorm.DB, opp Opportunity) error {
	result := conn.Save(&opp)
	if result.Error != nil {
		return fmt.Errorf("error saving opportunity %s: %w", opp.NoticeID, result.Error)
	}
	return nil
}

// SaveOpportunities is a helper to save a slice of opportunities
func SaveOpportunities(conn *gorm.DB, opps []Opportunity) error {
	for _, opp := range opps {
		err := SaveOpportunity(conn, opp)
		if err != nil {
			return fmt.Errorf("error saving opportunity %s: %w", opp.NoticeID, err)
		}
	}
	return nil
}
