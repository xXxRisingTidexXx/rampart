from typing import Optional
from requests import Session, codes
from requests.adapters import HTTPAdapter
from shapely.geometry import Point
from rampart.exceptions import RampartError
from rampart.models import Flat, Housing, Image, Interior


class Miner:
    __slots__ = ['_session', '_page', '_housings']

    def __init__(self, session: Session):
        self._session = session
        self._page = -1
        self._housings = {1: Housing.secondary, 2: Housing.primary}

    def mine_flat(self) -> Optional[Flat]:
        self._page += 1
        response = self._session.get(
            f'https://dom.ria.com/searchEngine/?category=1&realty_type=2&operation_type=1&full'
            f'CategoryOperation=1_2_1&limit=1&page={self._page}',
            timeout=30
        )
        if response.status_code != codes.ok:
            raise RampartError(f'Miner got status: {response.status_code}')
        items = response.json().get('items', [])
        if not len(items):
            self._page = -1
            return None
        url = items[0].get('beautiful_url', '')
        index = url.rindex('-')
        if index == -1:
            raise RampartError(f'Miner got URL without a dash: {url}')
        total_area = items[0].get('total_square_meters', 0.0)
        if total_area <= 0 or total_area > 1000:
            raise RampartError(f'Miner got invalid total area: {total_area}, {url}')
        living_area = items[0].get('living_square_meters', 0.0)
        if living_area < 0 or living_area > total_area:
            raise RampartError(f'Miner got invalid living area: {living_area}, {url}')
        kitchen_area = items[0].get('kitchen_square_meters', 0.0)
        if kitchen_area < 0 or kitchen_area > total_area:
            raise RampartError(f'Miner got invalid kitchen area: {kitchen_area}, {url}')
        room_number = items[0].get('rooms_count', 0)
        if room_number < 1 or room_number > 20:
            raise RampartError(f'Miner got invalid room number: {room_number}, {url}')
        total_floor = items[0].get('floors_count', 0)
        if total_floor < 1 or total_floor > 100:
            raise RampartError(f'Miner got invalid total floor: {total_floor}, {url}')
        floor = items[0].get('floor', 0)
        if floor < 1 or floor > total_floor:
            raise RampartError(f'Miner got invalid floor: {floor}, {url}')
        housing = items[0].get('realty_sale_type', -1)
        if housing not in self._housings:
            raise RampartError(f'Miner got invalid housing: {housing}, {url}')
        longitude = float(items[0].get('longitude', 0))
        if longitude < -180 or longitude > 180:
            raise RampartError(f'Miner got invalid longitude: {longitude}, {url}')
        latitude = float(items[0].get('latitude', 0))
        if latitude < -90 or latitude > 90:
            raise RampartError(f'Miner got invalid latitude: {latitude}, {url}')
        street = items[0].get('street_name_uk', '')
        if not street:
            street = items[0].get('street_name', '')
        return Flat(
            'https://dom.ria.com/uk/' + url,
            [
                Image(
                    f'https://cdn.riastatic.com/photosnew/dom/photo/{url[:index]}__{k}fl.webp',
                    Interior.unknown
                )
                for k in items[0].get('photos', {}).keys()
            ],
            float(items[0].get('priceArr', {}).get('1', '0').replace(' ', '')),
            total_area,
            living_area,
            kitchen_area,
            room_number,
            floor,
            total_floor,
            self._housings[housing],
            Point(longitude, latitude),
            items[0].get('city_name_uk', ''),
            street,
            items[0].get('building_number_str', ''),
            0,
            0,
            0
        )


def _main():
    with Session() as session:
        session.mount('https://', HTTPAdapter(pool_maxsize=5, max_retries=3))
        session.headers['User-Agent'] = 'RampartBot/1.0.0'
        print(Miner(session).mine_flat())


if __name__ == '__main__':
    _main()
