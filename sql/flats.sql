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

-- Reflect the closest to subway stations flats.
select origin_url, concat(st_y(point), ',', st_x(point)) as point, ssf
from flats
order by ssf desc
limit 30;

-- View the furthest to industrial zones flats.
select origin_url, concat(st_y(point), ',', st_x(point)) as point, izf
from flats
order by izf
limit 30;

-- Print the closest to industrial zones flats;
select origin_url, concat(st_y(point), ',', st_x(point)) as point, izf
from flats
order by izf desc
limit 30;

-- Discover the closest to parks apartments.
select origin_url, concat(st_y(point), ',', st_x(point)) as point, gzf
from flats
order by gzf desc
limit 30;
