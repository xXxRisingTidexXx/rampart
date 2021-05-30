create table if not exists images
(
    id       serial primary key not null,
    flat_id  integer            not null,
    url      varchar(256)       not null check ( url != '' ),
    interior smallint           not null check ( interior between 0 and 5 ),
    foreign key (flat_id) references flats (id),
    unique (flat_id, url)
);
