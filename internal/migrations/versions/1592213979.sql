create table if not exists flats
(
    id           serial primary key  not null,
    origin_url   varchar(500) unique not null,
    image_url    varchar(500)        not null,
    update_time  timestamp           not null,
    price        real                not null,
    total_area   real                not null,
    living_area  real                not null,
    kitchen_area real                not null,
    room_number  smallint            not null,
    floor        smallint            not null,
    total_floor  smallint            not null,
    housing      varchar(10)         not null,
    complex      varchar(50)         not null,
    point        geometry            not null,
    state        varchar(30)         not null,
    city         varchar(50)         not null,
    district     varchar(50)         not null,
    street       varchar(70)         not null,
    house_number varchar(10)         not null
);
