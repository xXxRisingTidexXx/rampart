create table if not exists recognitions
(
    id                  serial primary key not null,
    completion_time     timestamp          not null,
    abandoned_number    integer            not null check ( 0 <= abandoned_number ),
    luxury_number       integer            not null check ( 0 <= luxury_number ),
    comfort_number      integer            not null check ( 0 <= comfort_number ),
    junk_number         integer            not null check ( 0 <= junk_number ),
    construction_number integer            not null check ( 0 <= construction_number ),
    excess_number       integer            not null check ( 0 <= excess_number ),
    reading_duration    real               not null check ( 0 <= reading_duration ),
    loading_duration    real               not null check ( 0 <= loading_duration ),
    update_duration     real               not null check ( 0 <= update_duration ),
    total_duration      real               not null check ( 0 <= total_duration )
);