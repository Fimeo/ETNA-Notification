package repository

import (
	"net/http"

	"etna-notification/internal/infrastructure/etna"
)

type IEtnaRepository interface {
	Login(authentication etna.Authentication) (*http.Cookie, error)
	RetrieveNotifications(cookie *http.Cookie, username string) ([]etna.Notification, error)
}

func NewEtnaRepository() IEtnaRepository {
	return etna.Service{}
}
