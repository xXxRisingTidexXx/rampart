from typing import Optional
from requests import Session, codes
from shapely.geometry import Point
from rampart.exceptions import RampartError
from rampart.models import Flat


class Geocoder:
    __slots__ = ['_session']

    def __init__(self, session: Session):
        self._session = session

    def geocode_flat(self, flat: Flat) -> Optional[Flat]:
        if flat.has_location or not flat.city or not flat.street or not flat.house_number:
            return None
        response = self._session.get(
            f'https://nominatim.openstreetmap.org/search?format=json&countrycodes=ua&q='
            f'{flat.house_number.replace(" ", "+")}+{flat.street.replace(" ", "+")},+'
            f'{flat.city.replace(" ", "+")}',
            timeout=10
        )
        if response.status_code != codes.ok:
            raise RampartError(f'Geocoder got status: {response.status_code}')
        locations = response.json()
        if not len(locations):
            return None
        longitude = float(locations[0].get('lon', 0))
        if longitude < -180 or longitude > 180:
            raise RampartError(f'Geocoder got invalid longitude: {longitude}, {flat.url}')
        latitude = float(locations[0].get('lat', 0))
        if latitude < -90 or latitude > 90:
            raise RampartError(f'Geocoder got invalid latitude: {latitude}, {flat.url}')
        return Flat(
            flat.url,
            flat.images_urls,
            flat.price,
            flat.total_area,
            flat.living_area,
            flat.kitchen_area,
            flat.room_number,
            flat.floor,
            flat.total_floor,
            flat.housing,
            Point(longitude, latitude),
            flat.city,
            flat.street,
            flat.house_number,
            flat.ssf,
            flat.izf,
            flat.gzf
        )
