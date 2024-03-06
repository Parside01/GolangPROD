package main

import (
	"log"
	"log/slog"
	"os"
	"solution/server"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	logger := slog.Default()

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		logger.Error("missed SERVER_ADDRESS env (export smth like ':8080')")
		os.Exit(1)
	}
	server := server.NewServer(serverAddress, logger)

	server.Start()
}
