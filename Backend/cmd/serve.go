/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Server/internal/entities/cfg"
	"Server/internal/infrastructure"
	"Server/internal/server"
	"Server/internal/service"
	"Server/internal/storage"
	"Server/internal/storage/postgres"
	"Server/internal/storage/redis"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  `Starts the server & Listens for requests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing URL Shortener Backend...")

		// 1. Load Config
		appConfig := cfg.NewApp()
		_ = appConfig.Load()
		appConfig.Print()

		// 2. Initialize Infrastructure
		appInfra := infrastructure.NewApp()
		if err := appInfra.Connect(appConfig); err != nil {
			log.Fatalf("Failed to connect to infrastructure: %v", err)
		}

		// 3. Initialize Storage (normalized repos)
		userRepo := postgres.NewUserRepo(appInfra.PostgreSQL().DB)
		longURLRepo := postgres.NewLongURLRepo(appInfra.PostgreSQL().DB)
		shortURLRepo := postgres.NewShortURLRepo(appInfra.PostgreSQL().DB)
		statsRepo := postgres.NewStatsRepo(appInfra.PostgreSQL().DB)
		urlCache := redis.NewURLCache(appInfra.Redis().Client())

		appStorage := storage.NewStorage(userRepo, longURLRepo, shortURLRepo, statsRepo, urlCache)

		// 4. Initialize Services
		appService := service.NewService(appStorage, appConfig.Self.JWTSecret())

		// 5. Initialize Server
		appServer := server.NewServer(appService, appConfig.Self.JWTSecret())

		// 6. Start Serving
		port := fmt.Sprintf("%d", appConfig.Self.Port())
		if port == "0" {
			port = "8080"
		}

		if err := appServer.Serve(port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
