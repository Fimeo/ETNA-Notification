package domain

import "gorm.io/gorm"

type CalendarEvent struct {
	gorm.Model
	ExternalID int
	User       *User
	UserID     int
}
