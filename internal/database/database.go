package database

import (
	"gorm.io/gorm"
)

const (
	ConnectionUser     = "POSTGRES_USER"
	ConnectionPassword = "POSTGRES_PASSWORD"
	ConnectionHost     = "POSTGRES_HOST"
	ConnectionDatabase = "POSTGRES_DB"
	ConnectionPort     = "POSTGRES_PORT"
)

type Client struct {
	DB *gorm.DB
}
