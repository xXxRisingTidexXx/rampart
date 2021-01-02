create table if not exists subscriptions
(
    id          serial primary key not null,
    chat_id     bigint             not null,
    city        varchar(50)        not null check ( city != '' ),
    price       real               not null check ( 0 <= price ),
    room_number varchar(5)         not null check ( room_number != '' ),
    floor       varchar(5)         not null check ( floor != '' )
);
