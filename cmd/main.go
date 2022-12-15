package main

import (
	"go.uber.org/fx"

	"etna-notification/internal/config"
	"etna-notification/internal/controller"
	"etna-notification/internal/database"
	"etna-notification/internal/logging"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
	"etna-notification/pkg/security"
)

func main() {
	fx.New(
		fx.Provide(
			security.NewSecurity,
			service.InitServices,
			repository.InitRepositories,
			database.InitDatabaseConnection,
		),
		fx.Invoke(
			config.LoadConfig,
			logging.InitLogger,
			// repository.AutoMigrateModels,
		),
		fx.Invoke(controller.Setup),
	).Run()
}
