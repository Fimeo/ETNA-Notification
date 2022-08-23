package application

import (
	"os"

	"github.com/spf13/viper"

	"etna-notification/internal/application/repository"
	"etna-notification/internal/infrastructure/logger"
	"etna-notification/internal/infrastructure/mysql"
)

type Dependencies struct {
	Discord      repository.IDiscordRepository
	Notification repository.INotificationRepository
	Etna         repository.IEtnaRepository
	Users        repository.IUsersRepository
	f            *os.File
}

func (d Dependencies) Close() {
	d.Discord.Close()
	d.Notification.Close()
	d.f.Close()
}

func LoadDependencies() Dependencies {
	f := logger.Logger()

	dg, err := repository.NewDiscordRepository()
	if err != nil {
		panic(err)
	}

	database, err := mysql.NewConnection(viper.GetString("sqlite.file"))
	if err != nil {
		panic(err)
	}

	return Dependencies{
		Discord:      dg,
		Notification: database,
		Etna:         repository.NewEtnaRepository(),
		Users:        database,
		f:            f,
	}
}
