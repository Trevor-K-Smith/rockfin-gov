// Package database provides database-related commands
package database

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabaseCmd returns the database command group
func NewDatabaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "Database management commands",
	}

	cmd.AddCommand(newConnectionTestCmd())
	return cmd
}

func newConnectionTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connection-test",
		Short: "Test database connection using config settings",
		Run: func(cmd *cobra.Command, args []string) {
			host := viper.GetString("database.host")
			port := viper.GetInt("database.port")
			user := viper.GetString("database.user")
			password := viper.GetString("database.password")

			fmt.Printf("Attempting connection with:\n")
			fmt.Printf("Host: %s\n", host)
			fmt.Printf("Port: %d\n", port)
			fmt.Printf("Username: %s\n", user)
			fmt.Printf("Password: %s\n", password)

			dsn := fmt.Sprintf("host=%s port=%d dbname=postgres user=%s password=%s sslmode=disable",
				host, port, user, password)

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("Failed to get underlying *sql.DB: %v", err)
			}

			err = sqlDB.Ping()
			if err != nil {
				log.Fatalf("Failed to ping database: %v", err)
			}

			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("\n%s\n", green("Successfully connected to PostgreSQL server!"))
		},
	}
}
