package repository

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"

	"etna-notification/internal/infrastructure/discord"
)

type IDiscordRepository interface {
	SendTextMessage(channelID, message string) (*discordgo.Message, error)
	CreateUserNotificationTextChannel(username string) (*discordgo.Channel, error)
	GetChannel(channelID string) (*discordgo.Channel, error)
	OpenSocket() error
	Close()
}

func NewDiscordRepository() (IDiscordRepository, error) {
	dg, err := discord.NewConnection(viper.GetString("discord.bot-token"))
	if err != nil {
		return nil, err
	}

	err = dg.OpenSocket()
	if err != nil {
		return nil, err
	}

	return dg, nil
}
