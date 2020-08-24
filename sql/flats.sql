-- View housing and complex count stats.
select concat(housing, ' with complexes') as category, count(*) as flat_count
from flats
where complex != ''
group by housing
union
select concat(housing, ' without complexes') as category, count(*) as flat_count
from flats
where complex = ''
group by housing
union
select 'total' as category, count(*) as flat_count
from flats
order by category;

-- List all available hosing complexes.
select distinct complex
from flats
where complex != '';

-- Discover flat count by states.
select state, count(*) as flat_count
from flats
group by state
order by flat_count desc;

-- Discover flat count by cities.
select city, count(*) as flat_count
from flats
group by city
order by flat_count desc;

-- Discover complex count by cities.
select city, count(distinct complex) as complex_count
from flats
where complex != ''
group by city
order by complex_count desc;

-- Discover flat count by city districts.
select city, district, count(*) as flat_count
from flats
group by city, district
order by flat_count desc;

-- Explore cities with subway stations.
select city, count(*) as flat_count
from flats
where subway_station_distance != -1
group by city
order by flat_count desc;

-- View flats near subway from cities without subway
select origin_url, st_astext(point) as point
from flats
where city not in ('Київ', 'Харків', 'Дніпро')
  and subway_station_distance != -1;

select origin_url, concat(st_y(point), ',', st_x(point)) as point, subway_station_distance
from flats
where subway_station_distance != -1
order by subway_station_distance
limit 30;

select origin_url, concat(st_y(point), ',', st_x(point)) as point, industrial_zone_distance
from flats
where industrial_zone_distance != -1
order by industrial_zone_distance
limit 30;

select origin_url, concat(st_y(point), ',', st_x(point)) as point, industrial_zone_distance
from flats
where industrial_zone_distance != -1
order by industrial_zone_distance desc
limit 30;

select origin_url, concat(st_y(point), ',', st_x(point)) as point, green_zone_distance
from flats
where green_zone_distance != -1
order by green_zone_distance
limit 30;
