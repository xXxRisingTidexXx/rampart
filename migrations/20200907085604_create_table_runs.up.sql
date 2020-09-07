create table if not exists runs
(
    id                              serial primary key not null,
    completion_time                 timestamp          not null,
    miner                           varchar(40)        not null check ( miner != '' ),
    state_sanitation_number         smallint           not null check ( 0 <= state_sanitation_number ),
    city_sanitation_number          smallint           not null check ( 0 <= city_sanitation_number ),
    district_sanitation_number      smallint           not null check ( 0 <= district_sanitation_number ),
    swap_sanitation_number          smallint           not null check ( 0 <= swap_sanitation_number ),
    street_sanitation_number        smallint           not null check ( 0 <= street_sanitation_number ),
    house_number_sanitation_number  smallint           not null check ( 0 <= house_number_sanitation_number ),
    located_geocoding_number        smallint           not null check ( 0 <= located_geocoding_number ),
    unlocated_geocoding_number      smallint           not null check ( 0 <= unlocated_geocoding_number ),
    failed_geocoding_number         smallint           not null check ( 0 <= failed_geocoding_number ),
    inconclusive_geocoding_number   smallint           not null check ( 0 <= inconclusive_geocoding_number ),
    successful_geocoding_number     smallint           not null check ( 0 <= successful_geocoding_number ),
    subwayless_ssf_gauging_number   smallint           not null check ( 0 <= subwayless_ssf_gauging_number ),
    failed_ssf_gauging_number       smallint           not null check ( 0 <= failed_ssf_gauging_number ),
    inconclusive_ssf_gauging_number smallint           not null check ( 0 <= inconclusive_ssf_gauging_number ),
    successful_ssf_gauging_number   smallint           not null check ( 0 <= successful_ssf_gauging_number ),
    failed_izf_gauging_number       smallint           not null check ( 0 <= failed_izf_gauging_number ),
    inconclusive_izf_gauging_number smallint           not null check ( 0 <= inconclusive_izf_gauging_number ),
    successful_izf_gauging_number   smallint           not null check ( 0 <= successful_izf_gauging_number ),
    failed_gzf_gauging_number       smallint           not null check ( 0 <= failed_gzf_gauging_number ),
    inconclusive_gzf_gauging_number smallint           not null check ( 0 <= inconclusive_gzf_gauging_number ),
    successful_gzf_gauging_number   smallint           not null check ( 0 <= successful_gzf_gauging_number ),
    approved_validation_number      smallint           not null check ( 0 <= approved_validation_number ),
    uninformative_validation_number smallint           not null check ( 0 <= uninformative_validation_number ),
    denied_validation_number        smallint           not null check ( 0 <= denied_validation_number ),
    created_storing_number          smallint           not null check ( 0 <= created_storing_number ),
    updated_storing_number          smallint           not null check ( 0 <= updated_storing_number ),
    unaltered_storing_number        smallint           not null check ( 0 <= unaltered_storing_number ),
    failed_storing_number           smallint           not null check ( 0 <= failed_storing_number ),
    fetching_duration               real               not null check ( 0 <= fetching_duration ),
    geocoding_duration              real               not null check ( 0 <= geocoding_duration ),
    ssf_gauging_duration            real               not null check ( 0 <= ssf_gauging_duration ),
    izf_gauging_duration            real               not null check ( 0 <= izf_gauging_duration ),
    gzf_gauging_duration            real               not null check ( 0 <= gzf_gauging_duration ),
    reading_duration                real               not null check ( 0 <= reading_duration ),
    creation_duration               real               not null check ( 0 <= creation_duration ),
    update_duration                 real               not null check ( 0 <= update_duration ),
    total_duration                  real               not null check ( 0 <= total_duration )
);
