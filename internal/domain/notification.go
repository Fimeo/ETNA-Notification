package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	ExternalID int
	User       *User
	UserID     int
}
