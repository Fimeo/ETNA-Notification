package service

import (
	"github.com/imroc/req/v3"
	"go.uber.org/fx"

	"etna-notification/internal/logging"
)

type Service struct {
	IDiscordService
	IEtnaWebService
	logging.Logger
	client *req.Client
}

func InitServices(lc fx.Lifecycle) Service {
	client := NewClient()
	return Service{
		IDiscordService: NewDiscordService(lc),
		IEtnaWebService: NewEtnaWebservice(client),
		Logger:          logging.InitLogger(),
		client:          client,
	}
}
