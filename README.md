## ETNA Notification

Retrieve and send new notification from ETNA School Intranet to Discord Channel (bot).

This application could manage multiples users with their own discord private channel to
receive their notifications.

## Getting Started

### Prerequisites

* Docker
* A discord account and a fresh server
* A Discord Bot [Create a discord bot](https://github.com/Fimeo/ETNA-Notification/blob/main/doc/CreateDiscordBot.md)
* ETNA School account information

### Run the project

Run project with docker-compose : `docker-compose up -d`.

On the first launch, the database could be not ready at all when the app was launched.
Wait the postgres database server to be ready (see logs with `docker-compose logs -f postgres`) and then
restart migration and app services : `docker-compose restart migrate app`.

A cron was launch and retrieves notifications from intranet every 30 minutes. If notification was already send to your 
channel, no new notification was sent. Else a new discord message was written by the discord bot and allow you to stay 
informed from the latest news without perpetual connection.

Thanks to [DiscordGo](https://github.com/bwmarrin/discordgo)

---

This project is still in development:

TODOLIST :

* [ ] Register new users by a web interface
  * [ ] WEB UI to get account information
  * [ ] Create invitation link and rely discord account to etna user account
  * [ ] Create a new discord private channel for this user
* [ ] Fresh installation on a new server
  * [ ] Use some configuration to provide the Guild ID discord token
* [ ] Documentation
  * [ ] Create bot
  * [ ] Link to a server
