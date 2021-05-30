from typing import Optional
from urllib.parse import quote
from requests import Session, codes
from osm2geojson import xml2shapes
from rampart.exceptions import RampartError
from rampart.models import Flat


class Gauger:
    __slots__ = ['_session', '_host', '_subway_cities']

    def __init__(self, session: Session, host: str):
        self._session = session
        self._host = host
        self._subway_cities = {'Київ', 'Харків', 'Дніпро'}

    def gauge_flat(self, flat: Flat) -> Optional[Flat]:
        return None if not flat.has_location else Flat(
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
            flat.point,
            flat.city,
            flat.street,
            flat.house_number,
            self._gauge_ssf(flat),
            self._gauge_izf(flat),
            self._gauge_gzf(flat)
        )

    def _gauge_ssf(self, flat: Flat) -> float:
        if flat.city not in self._subway_cities:
            return 0
        _ = self._query_overpass(
            f'node[station=subway](around:2000,{flat.point.y},{flat.point.x});out;'
        )
        return 0.001

    def _query_overpass(self, query: str):
        response = self._session.get(
            f'https://{self._host}/api/interpreter?data={quote(query)}',
            timeout=30
        )
        if response.status_code != codes.ok:
            raise RampartError(f'Gauger got status: {response.status_code}')
        return xml2shapes(response.text)

    def _gauge_izf(self, flat: Flat) -> float:
        radius = 3000
        _ = self._query_overpass(
            f'(way[landuse=industrial](around:{radius},{flat.point.y},{flat.point.x});>;relation[l'
            f'anduse=industrial](around:{radius},{flat.point.y},{flat.point.x});>;);out;'
        )
        return 0

    def _gauge_gzf(self, flat: Flat) -> float:
        radius = 2500
        _ = self._query_overpass(
            f'(way[landuse=park](around:{radius},{flat.point.y},{flat.point.x});>;relation[landuse'
            f'=park](around:{radius},{flat.point.y},{flat.point.x});>;);out;'
        )
        return 0
