create table public.users
(
    id              bigserial
        primary key,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
    channel_id      text,
    discord_account text,
    login           text,
    password        text,
    status          text
);

create index idx_users_deleted_at
    on public.users (deleted_at);

create table public.notifications
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    external_id bigint,
    user_id     bigint
        constraint fk_notifications_user
            references public.users
);

create index idx_notifications_deleted_at
    on public.notifications (deleted_at);
