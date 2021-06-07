create table if not exists images
(
    id       serial primary key  not null,
    url      varchar(256) unique not null check ( url != '' ),
    interior varchar(15)         not null check ( interior != '' ) default 'unknown'
);
