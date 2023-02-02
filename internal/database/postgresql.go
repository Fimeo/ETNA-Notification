package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func InitDatabaseConnection() Client {
	connection, err := postgresConnection(
		os.Getenv(ConnectionUser),
		os.Getenv(ConnectionPassword),
		os.Getenv(ConnectionHost),
		os.Getenv(ConnectionPort),
		os.Getenv(ConnectionDatabase),
	)
	if err != nil {
		panic(err)
	}

	return Client{connection}
}

func postgresConnection(username, password, host, port, database string) (*gorm.DB, error) {
	dsn := "postgresql://" + username + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=disable"
	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[ERROR] Failed to connect to postgresql database : %+v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[ERROR]Cannot get sql db from gorm instance : %+v", err)
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Printf("[ERROR] Ping database failed : %+v", err)
		return nil, err
	}

	return db, nil
}
