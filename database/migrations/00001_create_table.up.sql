create table users
(
    id         bigserial
        primary key,
    user_id    bigint,
    time       timestamp with time zone,
    channel_id text,
    login      text,
    password   text
);

alter table users
    owner to postgres;

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

