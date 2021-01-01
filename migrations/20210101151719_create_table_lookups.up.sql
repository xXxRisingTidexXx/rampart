create table if not exists lookups
(
    id              serial primary key not null,
    subscription_id int                not null,
    status          varchar(10)        not null check ( status != '' ),
    foreign key (subscription_id) references subscriptions (id) on delete cascade
);
