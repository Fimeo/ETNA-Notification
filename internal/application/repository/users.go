package repository

import (
	"etna-notification/internal/infrastructure/mysql"
)

type IUsersRepository interface {
	GetEtnaUsers() ([]mysql.EtnaUser, error)
}
