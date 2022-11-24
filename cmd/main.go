package main

import (
	"log"

	"go.uber.org/fx"

	"github.com/joho/godotenv"

	"etna-notification/internal/controller"
	"etna-notification/internal/database"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

func main() {
	fx.New(
		fx.Invoke(loadConfig),
		fx.Invoke(service.NewLoggerService),
		fx.Provide(
			service.NewClient,
			service.NewDiscordService,
			service.NewEtnaWebservice,
			repository.NewUserRepository,
			repository.NewNotificationRepository,
			controller.NewEtnaNotificationController,
			database.NewDatabaseConnection,
		),
		fx.Invoke(controller.StartPushNotificationCron),
		fx.Invoke(controller.AutoMigrateModels),
	).Run()
	fx.OnStop(controller.CloseConnection)
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}