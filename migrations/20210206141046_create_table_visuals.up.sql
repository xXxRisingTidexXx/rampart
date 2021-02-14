create table if not exists visuals
(
    id       serial primary key not null,
    flat_id  integer            not null,
    image_id integer            not null,
    foreign key (flat_id) references flats (id) on delete cascade,
    foreign key (image_id) references images (id) on delete cascade,
    unique (flat_id, image_id)
);
