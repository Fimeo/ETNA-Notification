package main

import (
	"fmt"

	"github.com/robfig/cron"
	"github.com/spf13/viper"

	"etna-notification/internal/application"
	"etna-notification/internal/infrastructure/handler"
)

func main() {
	loadConfig()
	dependencies := application.LoadDependencies()
	defer dependencies.Close()

	cr := cron.New()
	cr.AddFunc("@every 30m", func() {
		handler.SendNewNotifications(dependencies)
	})
	cr.Start()
	handler.SendNewNotifications(dependencies)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
