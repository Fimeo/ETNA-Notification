## ETNA Notification

Retrieve and send new notification from ETNA School Intranet to Discord Channel (bot)

What you need :

* Sqlite3 installed
* Go installed
* Discord Bot

Connect to sqlite :  `sqlite3 notification.db`

Then run migrations : `cat migration/*.sql | sqlite3 database.db`

Copy and paste `config.override.yaml`

Run go project : `go run main.go`

A cron was launch and retrieves notifications from intranet every 15 minutes. If notification was already send to your channel, no new notification.
Else a new discord message was written by the discord bot and allow you to stay informed from the latest news without perpetual connection.