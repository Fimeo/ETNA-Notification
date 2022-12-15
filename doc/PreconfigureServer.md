# Preconfigure server

Add the discord bot in the fresh server.

---

Create a private `Notifications` category. Only bot and admin can view this category for the moment.

Copy the category id to fill `NOTIFICATION_CATEGORY_ID` configuration variable.

---

Create a private `System` category. Only bot and admin can view this category.

Then create a `error` channel. Copy the channel id and fill the `SYSTEM_ERROR_CHANNEL` configuration variable.

All errors are sent into this channel.

---

Create a public `connect` channel on general text channels of the server.

Copy the channel id and fill the `CONNECT_CHANNEL` configuration variable.

---

![CreateBot](https://github.com/Fimeo/ETNA-Notification/blob/main/doc/img/serverOverview.png)