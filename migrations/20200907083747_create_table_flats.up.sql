create table if not exists flats
(
    id           serial primary key    not null,
    url          varchar(256) unique   not null check ( url != '' ),
    price        real                  not null check ( 0 < price ),
    total_area   real                  not null check ( 0 < total_area ),
    living_area  real                  not null check ( 0 <= living_area and living_area <= total_area ),
    kitchen_area real                  not null check ( 0 <= kitchen_area and kitchen_area <= total_area ),
    room_number  smallint              not null check ( 0 < room_number ),
    floor        smallint              not null check ( 0 < floor and floor <= total_floor ),
    total_floor  smallint              not null check ( 0 < total_floor ),
    housing      smallint              not null check ( housing between 0 and 1),
    point        geometry(point, 4326) not null,
    city         varchar(50)           not null,
    street       varchar(70)           not null,
    house_number varchar(10)           not null,
    ssf          real                  not null,
    izf          real                  not null,
    gzf          real                  not null
);
