create table if not exists runs
(
    id serial primary key not null,
    completion_time timestamp not null,
    miner varchar(40) not null check ( miner != '' ),
    is_successful bool not null,
    located_geocoding_number smallint not null check ( 0 <= located_geocoding_number ),
    unlocated_geocoding_number smallint not null check ( 0 <= unlocated_geocoding_number ),
    failed_geocoding_number smallint not null check ( 0 <= failed_geocoding_number ),
    inconclusive_geocoding_number smallint not null check ( 0 <= inconclusive_geocoding_number ),
    successful_geocoding_number smallint not null check ( 0 <= successful_geocoding_number ),
    approved_validation_number smallint not null check ( 0 <= approved_validation_number ),
    denied_validation_number smallint not null check ( 0 <= denied_validation_number ),
    created_storing_number smallint not null check ( 0 <= created_storing_number ),
    updated_storing_number smallint not null check ( 0 <= updated_storing_number ),
    unaltered_storing_number smallint not null check ( 0 <= unaltered_storing_number ),
    failed_storing_number smallint not null check ( 0 <= failed_storing_number ),

    
);
