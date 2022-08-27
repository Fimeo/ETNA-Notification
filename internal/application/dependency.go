package application

import (
	"os"

	"etna-notification/internal/application/repository"
	"etna-notification/internal/infrastructure/database"
)

type Dependencies struct {
	Discord      repository.IDiscordRepository
	Notification repository.INotificationRepository
	Etna         repository.IEtnaRepository
	Users        repository.IUsersRepository
	db           *database.Service
}

func (d Dependencies) Close() {
	d.Discord.Close()
	d.db.Close()
}

func LoadDependencies() Dependencies {
	dg, err := repository.NewDiscordRepository()
	if err != nil {
		panic(err)
	}

	connection, err := database.NewPostgresConnection(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		panic(err)
	}

	return Dependencies{
		Discord:      dg,
		Notification: connection,
		Etna:         repository.NewEtnaRepository(),
		Users:        connection,
		db:           connection,
	}
}
