package repository

import (
	"etna-notification/internal/database"
	"etna-notification/pkg/security"
)

type Repositories struct {
	INotificationRepository
	IUserRepository
	ICalendarEventRepository
}

func InitRepositories(client database.Client, security security.Security) Repositories {
	return Repositories{
		NewNotificationRepository(client),
		NewUserRepository(client, security),
		NewCalendarEventRepository(client),
	}
}

func AutoMigrateModels(repositories Repositories) {
	repositories.INotificationRepository.Migrate()
	repositories.IUserRepository.Migrate()
	repositories.ICalendarEventRepository.Migrate()
}
