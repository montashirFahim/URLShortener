/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Server/internal/entities/cfg"
	"Server/internal/infrastructure"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database migrations",
	Long:  `Runs database migrations from internal/storage/postgres/migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migrations...")

		// 1. Load Config
		appConfig := cfg.NewApp()
		_ = appConfig.Load()

		// 2. Initialize Infrastructure
		appInfra := infrastructure.NewApp()
		if err := appInfra.Connect(appConfig); err != nil {
			log.Fatalf("Failed to connect to infrastructure: %v", err)
		}
		defer appInfra.PostgreSQL().Close()

		// 3. Find Migration Files
		migrationDir := "internal/storage/postgres/migrations"
		files, err := os.ReadDir(migrationDir)
		if err != nil {
			log.Fatalf("Failed to read migration directory: %v", err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), "_up.sql") {
				fmt.Printf("Applying migration: %s\n", file.Name())
				content, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
				if err != nil {
					log.Fatalf("Failed to read migration file %s: %v", file.Name(), err)
				}

				// Basic SQL parser: split by ; and handle comments correctly
				lines := strings.Split(string(content), "\n")
				var cleanSQL strings.Builder
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" || strings.HasPrefix(trimmed, "--") {
						continue
					}
					cleanSQL.WriteString(line + "\n")
				}

				queries := strings.Split(cleanSQL.String(), ";")
				for _, query := range queries {
					q := strings.TrimSpace(query)
					if q == "" {
						continue
					}
					_, err := appInfra.PostgreSQL().DB.Exec(q)
					if err != nil {
						log.Fatalf("Failed to execute query in %s: %v\nQuery: %s", file.Name(), err, q)
					}
				}
			}
		}
		fmt.Println("Migrations completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
