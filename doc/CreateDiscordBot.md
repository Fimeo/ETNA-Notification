# How to create a discord Bot

First, you need to create a Discord account.

**Enable developer mode** :

Click on User Setting button to enable developer mode in your Discord account.

![Settings](https://github.com/Fimeo/ETNA-Notification/blob/main/doc/img/settings.png)

Then, click on Advanced and then enable Developer Mode.

**Create a Discord application** :

Go to https://discord.com/developers/applications/, then create a new application with "New Application" button, and then name the app and click on "Create" button.

**Create a Bot** :

Click on Bot menu and then on Add Bot button.

![CreateBot](https://github.com/Fimeo/ETNA-Notification/blob/main/doc/img/createBot.png)


Now go in OAuth2 menu and click on Copy button in order to get Client ID information.

**Generate the Bot invite link** :

In order to link our Bot to one of our Discord server, we need to generate an invite link, with the CLIENT ID we copied:
```https://discord.com/api/oauth2/authorize?client_id=<CLIENT-ID>&permissions=8&scope=bot```

When we go to this URL, a connection window appears to link the bot to a discord server.

**Save the token** :

There is one last thing to do so that our Go application can connect to the Discord server: we need a token.

For that, go back in the Discord application developers website, then click on Bot menu and then click on Copy button in order to copy the token (and save it somewhere).

It's time for the show !
