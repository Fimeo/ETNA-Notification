package domain

import (
	"time"
)

type EtnaAuthenticationInput struct {
	ID             int      `json:"id"`
	Login          string   `json:"login"`
	Email          string   `json:"email"`
	LogAs          bool     `json:"logas"`
	Groups         []string `json:"groups"`
	LoginDate      string   `json:"login_date"`
	FirstConnexion bool     `json:"firstconnexion"`
	Password       string   `json:"password"`
}

type EtnaNotification struct {
	ID                    int         `json:"id"`
	Message               string      `json:"message"`
	Start                 time.Time   `json:"start"`
	End                   interface{} `json:"end"`
	CanValidate           bool        `json:"can_validate"`
	Validated             bool        `json:"validated"`
	Type                  string      `json:"type"`
	EtnaNotificationMetas `json:"metas"`
}

type EtnaNotificationMetas struct {
	Type         string `json:"type"`
	SessionID    int    `json:"session_id,omitempty"`
	ActivityType string `json:"activity_type,omitempty"`
	ActivityID   int    `json:"activity_id,omitempty"`
	Promo        string `json:"promo,omitempty"`
}

func BuildAuthenticationFromUser(user *User) *EtnaAuthenticationInput {
	return &EtnaAuthenticationInput{
		ID:             user.UserID,
		Login:          user.Login,
		Email:          user.Login + "@etna-alternance.net",
		LogAs:          false,
		Groups:         []string{"student"},
		LoginDate:      time.Now().Format("2006-01-02 15-04-05"),
		FirstConnexion: false,
		Password:       user.Password,
	}
}

func BuildNotificationFromEtnaNotificationAndUser(notification *EtnaNotification, user *User) *Notification {
	return &Notification{
		ExternalID: notification.ID,
		User:       user,
	}
}

func BuildMessageFromEtnaNotification(notification *EtnaNotification) string {
	return ":bell: " + "[" + notification.Type + "]" + notification.Message
}
