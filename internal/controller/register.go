package controller

import (
	"github.com/bwmarrin/discordgo"
	"os"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
	"etna-notification/internal/usecase"
)

type registerController struct {
	DiscordService         service.IDiscordService
	EtnaWebService         service.IEtnaWebService
	UserRepository         repository.IUserRepository
	NotificationRepository repository.INotificationRepository
}

type IRegisterController interface {
	Register(login, password, discordAccount string) (*string, error)
	Connect()
}

func NewRegisterController(repositories repository.Repositories, services service.Service) IRegisterController {
	return &registerController{
		DiscordService:         services.IDiscordService,
		EtnaWebService:         services.IEtnaWebService,
		UserRepository:         repositories.IUserRepository,
		NotificationRepository: repositories.INotificationRepository,
	}
}

// Register takes account information, try to authenticate ETNA Account and save credentials.
// If user login is already registered, update the password and discord account fields only. Returns the invitation
// link to join the server.
func (c *registerController) Register(login, password, discordAccount string) (*string, error) {
	user := &domain.User{
		DiscordAccount: discordAccount,
		Login:          login,
		Password:       password,
	}

	err := usecase.CheckEtnaAccountValidity(c.EtnaWebService, login, password)
	if err != nil {
		return nil, err
	}

	err = usecase.RegisterNewUser(c.UserRepository, user)
	if err != nil {
		return nil, err
	}

	link, err := usecase.CreateServerInvitation(c.DiscordService, os.Getenv(service.ConnectChannelID))
	if err != nil {
		return nil, err
	}

	return link, nil
}

// Connect register connect handler func that follow message creation in connect channel. The connect method checks
// account creation and create a personal channel for the user to receive their notifications.
func (c *registerController) Connect() {
	c.DiscordService.Session().AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		err := usecase.CreatePersonalChannel(c.UserRepository, c.DiscordService, m)
		if err != nil {
			usecase.SendErrorNotification(c.DiscordService, err.Error())
		}
	})
}
