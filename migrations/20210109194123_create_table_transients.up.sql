create table if not exists transients
(
    id          bigint primary key not null,
    status      varchar(15)        not null check ( status != '' )      default 'city',
    city        varchar(50)        not null check ( city != '' ),
    price       real               not null check ( 0 <= price )        default 0,
    room_number varchar(10)        not null check ( room_number != '' ) default 'any',
    floor       varchar(10)        not null check ( floor != '' )       default 'any'
);
