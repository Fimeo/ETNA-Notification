create table users
(
    id              bigserial
        primary key,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
    time            timestamp with time zone,
    channel_id      text,
    discord_account text,
    login           text,
    password        text,
    status          text
);

alter table users
    owner to postgres;

create index idx_users_deleted_at
    on users (deleted_at);

create table notifications
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    external_id bigint,
    user_id     bigint
        constraint fk_notifications_user
            references users
);

alter table notifications
    owner to postgres;

create index idx_notifications_deleted_at
    on notifications (deleted_at);
