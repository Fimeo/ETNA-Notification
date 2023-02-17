create table public.calendar_events
(
    id          bigserial
        primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone,
    external_id bigint,
    user_id     bigint
        constraint fk_calendar_events_user
            references public.users
);

create index idx_calendar_events_deleted_at
    on public.calendar_events (deleted_at);

create index idx_calendar_events_external_id
    on public.calendar_events (external_id);

create index idx_notification_external_id
    on public.notifications (external_id);
