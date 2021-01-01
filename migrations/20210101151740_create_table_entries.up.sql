create table if not exists entries
(
    id        serial primary key not null,
    lookup_id integer            not null,
    flat_id   integer            not null,
    position  smallint           not null check ( 0 <= position ),
    foreign key (lookup_id) references lookups (id) on delete cascade,
    foreign key (flat_id) references flats (id) on delete cascade,
    unique (lookup_id, flat_id)
);
