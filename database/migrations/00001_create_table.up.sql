CREATE TABLE notification
(
    id          INTEGER NOT NULL PRIMARY KEY,
    time        DATE NOT NULL,
    external_id INTEGER NOT NULL,
    "user"      VARCHAR NOT NULL
);

CREATE TABLE users
(
    id       INTEGER NOT NULL PRIMARY KEY,
    time     DATE NOT NULL,
    user_id  INTEGER NOT NULL,
    channelID VARCHAR NOT NULL,
    login    VARCHAR NOT NULL,
    password VARCHAR NOT NULL
);
