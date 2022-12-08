package service

import (
	"context"
	"go.uber.org/fx"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"etna-notification/internal/domain"
)

// DiscordService is a proxy using discord-go package.
type discordService struct {
	DG *discordgo.Session
}

type IDiscordService interface {
	GetCurrentSession() *discordgo.Session
	SendTextMessage(channelID, message string) (*discordgo.Message, error)
	CreateUserNotificationTextChannel(username *domain.User) (*discordgo.Channel, error)
	ChannelNewReadingMember(memberID, channelID string) (*discordgo.Channel, error)
	GetChannel(channelID string) (*discordgo.Channel, error)
	CreateInvitation() (*discordgo.Invite, error)
	CloseConnection()
}

// NewDiscordService init connection with discord using discord-go package. A web socket is open.
// Use the fx lifecycle to close the web socket connection
func NewDiscordService(lc fx.Lifecycle) IDiscordService {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Panicf("[ERROR] Failed to init discord connection using token : %+v", err)
	}

	err = dg.Open()
	if err != nil {
		log.Panicf("[ERROR] Failed to open discord web socket : %+v", err)
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return dg.Close()
		},
	})

	return &discordService{DG: dg}
}

func (dg *discordService) GetCurrentSession() *discordgo.Session {
	return dg.DG
}

func (dg *discordService) SendTextMessage(channelID, message string) (*discordgo.Message, error) {
	messageSend, err := dg.DG.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("[ERROR] Failed to send discord message : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
	return messageSend, nil
}

func (dg *discordService) CreateUserNotificationTextChannel(user *domain.User) (*discordgo.Channel, error) {
	guildID := "1050391937783435284" // The server guild ID // TODO : use configuration
	channelCreate, err := dg.DG.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     user.Login,
		Type:     discordgo.ChannelTypeGuildText,
		Position: 1,
		ParentID: "1050433185932116099", // The category Notification ID // TODO : use configuration
	})
	if err != nil {
		log.Printf("[ERROR] Failed to create guild channel : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] Guild channel created : ", channelCreate.Name)
	return channelCreate, nil
}

func (dg *discordService) GetChannel(channelID string) (*discordgo.Channel, error) {
	channel, err := dg.DG.Channel(channelID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve guid channel : %+v", err)
		return nil, err
	}

	return channel, nil
}

func (dg *discordService) ChannelNewReadingMember(memberID, channelID string) (*discordgo.Channel, error) {
	channel, err := dg.GetChannel(channelID)
	if err != nil {
		return nil, err
	}
	members := channel.PermissionOverwrites // Retrieves default members from override category
	members = append(members, &discordgo.PermissionOverwrite{
		ID:    memberID, // Add the new user notification account
		Type:  discordgo.PermissionOverwriteTypeMember,
		Deny:  0,
		Allow: 1024, // Reading only
	})

	channelUpdated, err := dg.DG.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
		PermissionOverwrites: members,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to add new channel member : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] New member added to channel : ", channelUpdated.Name)
	return channelUpdated, nil
}

func (dg *discordService) CreateInvitation() (*discordgo.Invite, error) {
	invitation, err := dg.DG.ChannelInviteCreate("1050392032079786077", discordgo.Invite{ // TODO : use configuration
		MaxAge:    86400, // 1 day
		MaxUses:   1,
		Temporary: false,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to create invitation : %+v", err)
		return nil, err
	}
	log.Print("[DEBUG] New invitation created : ", invitation.CreatedAt)
	return invitation, nil
}

func (dg *discordService) CloseConnection() {
	dg.DG.Close()
}
