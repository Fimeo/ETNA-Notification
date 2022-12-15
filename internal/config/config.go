package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"etna-notification/internal/database"
	"etna-notification/internal/service"
	"etna-notification/pkg/security"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mandatoryKey := []string{
		service.ConnectChannelID,
		service.NotificationCategoryID,
		service.SystemErrorChannelID,
		service.DumpRequest,
		database.ConnectionDatabase,
		database.ConnectionHost,
		database.ConnectionPassword,
		database.ConnectionPort,
		database.ConnectionUser,
		security.RsaPublicKey,
		security.RsaPrivateKey,
	}
	for _, s := range mandatoryKey {
		if os.Getenv(s) == "" {
			fmt.Print(s + " env key is missing")
			os.Exit(1)
		}
	}
}
