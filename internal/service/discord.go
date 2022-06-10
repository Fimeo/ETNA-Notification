package service

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	Session *discordgo.Session
}

func DiscordConn(bottoken string) DiscordService {
	s, err := discordgo.New("Bot " + bottoken)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to discord bot : %+v", err))
	}
	return DiscordService{Session: s}
}

func (s *DiscordService) SendDiscordMessage(channel, message string) {
	messageSend, err := s.Session.ChannelMessageSend(channel, message)
	if err != nil {
		return
	}
	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
}
