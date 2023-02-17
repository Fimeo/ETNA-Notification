package domain

import (
	"fmt"
	"time"
)

const (
	Notice = "notice"
	Error  = "error"
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

type CalendarEvent struct {
	ID    int    `json:"id"`
	Event int    `json:"event"`
	Name  string `json:"name"`
	// The slug of calendar event
	ActivityName string `json:"activity_name"`
	SessionName  string `json:"session_name"`
	// Type values: presential, suivi, soutenance
	Type     string `json:"type"`
	Location string `json:"location"`
	// Start time of the calendar event
	Start string `json:"start"`
	// End time of the calendar event
	End string `json:"end"`
	// Members concerned by the event
	Group        CalendarEventGroup        `json:"group"`
	Registration CalendarEventRegistration `json:"registration"`
	// UvName module name
	UvName string `json:"uv_name"`
}

type CalendarEventMember struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Validation int    `json:"validation"`
	Forced     int    `json:"forced"`
}

type CalendarEventRegistration struct {
	Type   string `json:"type"`
	Date   string `json:"date"`
	Forced int    `json:"forced"`
	Locked int    `json:"locked"`
}

func (e CalendarEvent) BuildCalendarMessage() string {
	return fmt.Sprintf(
		"%s : %s. %s : %s - %s",
		e.Name, e.ActivityName, e.Location, e.Start, e.End)
}

type CalendarEventGroup struct {
	ID         int                   `json:"id"`
	Leader     CalendarEventMember   `json:"leader"`
	Validation interface{}           `json:"validation"`
	Members    []CalendarEventMember `json:"members"`
}

func BuildNotificationFromEtnaNotificationAndUser(notification *EtnaNotification, user *User) *Notification {
	return &Notification{
		ExternalID: notification.ID,
		UserID:     int(user.ID),
	}
}

func BuildMessageFromEtnaNotification(notification *EtnaNotification) string {
	if notification.Type == Notice {
		return ":bell: " + notification.Message
	} else if notification.Type == Error {
		return ":x: " + notification.Message
	}
	return ":pushpin: " + notification.Message
}
