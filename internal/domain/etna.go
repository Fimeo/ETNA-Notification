package domain

import (
	"time"
)

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

func BuildNotificationFromEtnaNotificationAndUser(notification *EtnaNotification, user *User) *Notification {
	return &Notification{
		ExternalID: notification.ID,
		User:       user,
	}
}

func BuildMessageFromEtnaNotification(notification *EtnaNotification) string {
	return ":bell: " + "[" + notification.Type + "]" + notification.Message
}
