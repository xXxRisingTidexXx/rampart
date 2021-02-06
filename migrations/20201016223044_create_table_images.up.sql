create table if not exists images
(
    id       serial primary key not null,
    flat_id  integer            not null,
    url      varchar(256)       not null check ( url != '' ),
    kind     varchar(10)        not null check ( kind != '' ),
    interior varchar(15)        not null check ( interior != '' ) default 'unknown',
    foreign key (flat_id) references flats (id) on delete cascade,
    unique (flat_id, url)
);
