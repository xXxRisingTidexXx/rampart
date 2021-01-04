create table if not exists lookups
(
    id              serial primary key not null,
    subscription_id integer            not null,
    flat_id         integer            not null,
    status          varchar(10)        not null check ( status != '' ),
    foreign key (subscription_id) references subscriptions (id) on delete cascade,
    foreign key (flat_id) references flats (id) on delete cascade,
    unique (subscription_id, flat_id)
);
