package domain

import "time"

type User struct {
	ID        int
	UserID    int
	Time      time.Time
	ChannelID string
	Login     string
	Password  string
}
