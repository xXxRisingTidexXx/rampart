create table if not exists transients
(
    id          bigint primary key not null,
    status      varchar(10)        not null check ( status != '' ),
    city        varchar(50)        not null check ( city != '' ),
    price       real               not null check ( 0 <= price ),
    room_number varchar(10)        not null check ( room_number != '' ),
    floor       varchar(10)        not null check ( floor != '' )
);
