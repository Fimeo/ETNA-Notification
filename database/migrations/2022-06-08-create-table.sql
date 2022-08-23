CREATE TABLE notification
(
    id          INTEGER NOT NULL PRIMARY KEY,
    time        DATETIME NOT NULL,
    external_id INTEGER NOT NULL,
    "user"      TEXT NOT NULL
);

CREATE TABLE users
(
    id       INTEGER NOT NULL PRIMARY KEY,
    time     DATETIME NOT NULL,
    user_id  INTEGER NOT NULL,
    channelID TEXT NOT NULL,
    login    TEXT NOT NULL,
    password TEXT NOT NULL
);