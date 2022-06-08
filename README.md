## ETNA Notification

Retrieve and send new notification from ETNA School Intranet to Discord Channel (bot)

What you need :

* Sqlite3 installed
* Go installed
* Discord Bot

Connect to sqlite :  `sqlite3 notification.db`

Then run migrations : `cat migration/*.sql | sqlite3 database.db`

Run go project : `go run main.go`