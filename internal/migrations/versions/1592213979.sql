create table if not exists flats
(
    id           serial primary key    not null,
    origin_url   varchar(500) unique   not null check ( origin_url != '' ),
    image_url    varchar(500)          not null,
    update_time  timestamp             not null,
    parsing_time timestamp             not null,
    price        real                  not null check ( 0 < price ),
    total_area   real                  not null check ( 0 < total_area ),
    living_area  real                  not null check ( 0 <= living_area and living_area < total_area ),
    kitchen_area real                  not null check ( 0 <= kitchen_area and kitchen_area < total_area ),
    room_number  smallint              not null check ( 0 < room_number ),
    floor        smallint              not null check ( 0 < floor and floor <= total_floor ),
    total_floor  smallint              not null check ( 0 < total_floor ),
    housing      varchar(10)           not null check ( housing != '' ),
    complex      varchar(50)           not null,
    point        geometry(point, 4326) not null,
    state        varchar(30)           not null,
    city         varchar(50)           not null,
    district     varchar(50)           not null,
    street       varchar(70)           not null,
    house_number varchar(20)           not null
);
