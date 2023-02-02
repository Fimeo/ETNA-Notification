package service

import (
	"context"
	"go.uber.org/fx"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"

	"etna-notification/internal/domain"
)

const (
	NotificationCategoryID = "NOTIFICATION_CATEGORY_ID"
	SystemErrorChannelID   = "SYSTEM_ERROR_CHANNEL"
	ConnectChannelID       = "CONNECT_CHANNEL"
)

// DiscordService is a proxy using discord-go package.
type discordService struct {
	DG *discordgo.Session
}

type IDiscordService interface {
	Session() *discordgo.Session
	SendTextMessage(channelID, message string) (*discordgo.Message, error)
	SendTextMessageReply(content string, message *discordgo.Message)
	CreateUserNotificationTextChannel(username *domain.User, guildID string) (*discordgo.Channel, error)
	ChannelNewReadingMember(memberID, channelID string) (*discordgo.Channel, error)
	GetChannel(channelID string) (*discordgo.Channel, error)
	CreateInvitation(channelID string) (*discordgo.Invite, error)
	ReplyInteractiveCommand(content string, i *discordgo.InteractionCreate)
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

// Session returns the discordgo.Session
func (dg *discordService) Session() *discordgo.Session {
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

func (dg *discordService) CreateUserNotificationTextChannel(user *domain.User, guildID string) (*discordgo.Channel, error) {
	channelCreate, err := dg.DG.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     user.Login,
		Type:     discordgo.ChannelTypeGuildText,
		Position: 1,
		ParentID: os.Getenv(NotificationCategoryID), // The category Notification ID
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

func (dg *discordService) CreateInvitation(channelID string) (*discordgo.Invite, error) {
	invitation, err := dg.DG.ChannelInviteCreate(channelID, discordgo.Invite{
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

func (dg *discordService) SendTextMessageReply(content string, message *discordgo.Message) {
	_, err := dg.DG.ChannelMessageSendReply(message.ChannelID, content, &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	})
	if err != nil {
		log.Printf("[ERROR] message reply has failed with content : %s %+v %+v", content, message, err)
		return
	}
}

func (dg *discordService) ReplyInteractiveCommand(content string, i *discordgo.InteractionCreate) {
	err := dg.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Printf("[ERROR] message interaction reply has failed with content : %s %+v", content, err)
		return
	}
}

func (dg *discordService) CloseConnection() {
	dg.DG.Close()
}
