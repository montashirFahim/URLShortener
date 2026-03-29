/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"User/internal/domain/cfg"
	"User/internal/infrastructure"
	"User/internal/repository/postgres"
	"User/internal/server/http"
	"User/internal/service"
	"log"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "manges user activities",
	Long:  `manges user created short urls`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load Config
		cnf := cfg.NewApp()
		if err := cnf.Load(); err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		// 2. Initialize Infrastructure
		app := infrastructure.NewApp()
		if err := app.Connect(cnf); err != nil {
			log.Fatalf("failed to connect to infra: %v", err)
		}

		// 3. Initialize Repositories
		userRepo := postgres.NewPostgresRepo(app.PostgreSQL().DB)
		urlRepo := postgres.NewUrlPostgresRepo(app.PostgreSQL().DB)

		// 4. Initialize Services
		userSvc := service.NewUserService(userRepo, []byte(cnf.JWTSecret.Secret.String()))
		urlSvc := service.NewUrlService(urlRepo)
		svc := service.NewService(userSvc, urlSvc)

		// 5. Initialize & Start Server
		srv := http.NewServer(app.PostgreSQL().DB, nil, svc, []byte(cnf.JWTSecret.Secret.String()))
		if err := srv.Serve(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
