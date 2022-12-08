## How to invite new people to register and receive notification on discord

### Workflow :

The new user goes to the webpage to register.
The website indicates that we need the password and this tool is not official. If you does not trust this service, you
can host him yourself using the public repo.

The user need to input some informations :
- username or email
- password

Back :

- Try connection to intranet, retrieve personal information (account id)

Next, the new member enter his discord name with hashtag

Next, if login succeed, we need to create discord invitation link

The invitation was display to user

When the user accept the invitation, he was bring onto the server in /slash channel.

Typing /connect will validate the account previously entered by the new member and create the channel, insert the account and trigger notifications (or just tell with a message is ok)

(Do notifications are saved in database directly for existing one or 1000+ notifications displayed on the first time ?)


---

usecase user has the invitation link but didnt register : send website link and abort /connect method.



## Other things to do :

Handle incorrect password if he change

Send status notifications into dedicated channel (intra down, service down)

Send errors on private channel (only admin)

