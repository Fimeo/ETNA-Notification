alter table calendar_events
    alter column external_id type varchar using external_id::varchar;
