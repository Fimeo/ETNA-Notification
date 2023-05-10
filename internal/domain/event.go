package domain

import "gorm.io/gorm"

type CalendarEvent struct {
	gorm.Model
	ExternalID string
	User       *User
	UserID     int
}
