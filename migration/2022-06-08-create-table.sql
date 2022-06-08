create table notification
(
    id          INTEGER  not null
        primary key,
    time        DATETIME not null,
    external_id INTEGER  not null,
    user        TEXT     not null
);