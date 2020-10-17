create table if not exists images
(
    id           serial primary key  not null,
    flat_id      int                 not null,
    url          varchar(256) unique not null check ( url != '' ),
    kind         varchar(10)         not null check ( kind != '' ),
    parsing_time timestamp           not null,
    constraint fk_images_flats foreign key (flat_id) references flats (id) on delete cascade
);
