package model

import "time"

type Authentication struct {
	ID             int      `json:"id"`
	Login          string   `json:"login"`
	Email          string   `json:"email"`
	Logas          bool     `json:"logas"`
	Groups         []string `json:"groups"`
	LoginDate      string   `json:"login_date"`
	Firstconnexion bool     `json:"firstconnexion"`
	Password       string   `json:"password"`
}

type Notification struct {
	ID          int         `json:"id"`
	Message     string      `json:"message"`
	Start       time.Time   `json:"start"`
	End         interface{} `json:"end"`
	CanValidate bool        `json:"can_validate"`
	Validated   bool        `json:"validated"`
	Type        string      `json:"type"`
	Metas       `json:"metas"`
}

type Metas struct {
	Type         string `json:"type"`
	SessionID    int    `json:"session_id,omitempty"`
	ActivityType string `json:"activity_type,omitempty"`
	ActivityID   int    `json:"activity_id,omitempty"`
	Promo        string `json:"promo,omitempty"`
}
