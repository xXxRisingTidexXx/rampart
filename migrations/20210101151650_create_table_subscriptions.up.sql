create table if not exists subscriptions
(
    id          serial primary key not null,
    uuid        uuid unique        not null,
    chat_id     bigint             not null,
    status      varchar(10)        not null check ( status != '' ) default 'actual',
    city        varchar(50)        not null check ( city != '' ),
    price       real               not null check ( 0 <= price ),
    room_number varchar(10)        not null check ( room_number != '' ),
    floor       varchar(10)        not null check ( floor != '' )
);
