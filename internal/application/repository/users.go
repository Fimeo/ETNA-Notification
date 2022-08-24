package repository

import (
	"etna-notification/internal/infrastructure/database"
)

type IUsersRepository interface {
	GetEtnaUsers() ([]database.EtnaUser, error)
}
