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
select city, count(complex) as complex_count
from flats
where complex != ''
group by city
order by complex_count desc;

-- Discover flat count by city districts.
select city, district, count(*) as flat_count
from flats
group by city, district
order by flat_count desc;
