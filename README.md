## ETNA Notification

Retrieve and send new notification from ETNA School Intranet to Discord Channel (bot)

What you need :

* Sqlite3 installed
* Go installed
* A Discord Bot [Create a discord bot](https://github.com/Fimeo/ETNA-Notification/doc/CreateDiscordBot.md)
* ETNA School account information

Connect to sqlite :  `sqlite3 notification.db`

Then run migrations : `cat migration/*.sql | sqlite3 database.db`

Copy and paste `config.override.yaml` as `config.yaml`

Run go project : `go run main.go`

A cron was launch and retrieves notifications from intranet every 30 minutes. If notification was already send to your channel, no new notification.
Else a new discord message was written by the discord bot and allow you to stay informed from the latest news without perpetual connection.

Thanks to https://github.com/bwmarrin/discordgo
