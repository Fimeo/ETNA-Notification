package database

import (
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}
