package domain

import (
	"fmt"
	"log"
	"time"
)

const (
	Notice         = "notice"
	Error          = "error"
	TypeSuivi      = "suivi"
	TypeSoutenance = "soutenance"
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

type EtnaCalendarEvent struct {
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
	Group        EtnaCalendarEventGroup        `json:"group"`
	Registration EtnaCalendarEventRegistration `json:"registration"`
	// UvName module name
	UvName string `json:"uv_name"`
}

type EtnaCalendarEventMember struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Validation int    `json:"validation"`
	Forced     int    `json:"forced"`
}

type EtnaCalendarEventRegistration struct {
	Type   string `json:"type"`
	Date   string `json:"date"`
	Forced int    `json:"forced"`
	Locked int    `json:"locked"`
}

type EtnaCalendarEventGroup struct {
	ID         int                       `json:"id"`
	Leader     EtnaCalendarEventMember   `json:"leader"`
	Validation interface{}               `json:"validation"`
	Members    []EtnaCalendarEventMember `json:"members"`
}

// IsNotifiable Returns true is the event type request is relevant.
func (e EtnaCalendarEvent) IsNotifiable() bool {
	return e.Type == TypeSuivi || e.Type == TypeSoutenance
}

// IsInNext30Minutes returns true is the event start date is between current time and current time + 30 minutes.
func (e EtnaCalendarEvent) IsInNext30Minutes() bool {
	currentTime := time.Now()
	eventStart, err := time.Parse("2006-01-02 15:04:05", e.Start)
	if err != nil {
		log.Printf("[ERROR] cannot parse input event start date : %s", e.Start)
		return false
	}
	difference := currentTime.Sub(eventStart)
	if difference.Minutes() > 0 && difference.Minutes() <= 30 {
		return true
	}
	return false
}

func (e EtnaCalendarEvent) BuildCalendarMessage() string {
	return fmt.Sprintf(
		":date: **%s** %s : %s. %s : %s - %s",
		e.UvName, e.Name, e.ActivityName, e.Location, e.Start, e.End)
}

func (n EtnaNotification) BuildNotificationMessage() string {
	switch n.Type {
	case Notice:
		return ":bell: " + n.Message
	case Error:
		return ":x: " + n.Message
	default:
		return ":pushpin: " + n.Message
	}
}

func (n EtnaNotification) BuildNotification(user *User) *Notification {
	return &Notification{
		ExternalID: n.ID,
		UserID:     int(user.ID),
	}
}

func (e EtnaCalendarEvent) BuildCalendarEvent(user *User) *CalendarEvent {
	return &CalendarEvent{
		ExternalID: e.ID,
		UserID:     int(user.ID),
	}
}
