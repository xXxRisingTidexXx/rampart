from os import getenv
from typing import List
from sqlalchemy import create_engine
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine

_engine: Engine = create_engine(getenv('RAMPART_DATABASE_DSN'))
_booster = Booster(model_file='model.txt')


def search_flats(query: 'Query') -> List['Flat']:
    frame = _read_flats(query)
    frame['score'] = _booster.predict(
        frame.drop(columns=['url', 'street', 'house_number']),
        num_iteration=_booster.best_iteration
    )
    return [
        Flat(
            s['url'],
            query.city,
            s['street'],
            s['house_number'],
            s['actual_price'],
            s['total_area'],
            s['actual_room_number'],
            s['actual_floor'],
            s['total_floor']
        )
        for _, s
        in frame.sort_values('score', ascending=False).head(3).iterrows()
    ]


class Query:
    __slots__ = ['city', 'price', 'floor', 'room_number']

    def __init__(self, city: str, price: float, floor: int, room_number: int):
        self.city = city
        self.price = price
        self.floor = floor
        self.room_number = room_number


class Flat:
    __slots__ = [
        'url',
        'city',
        'street',
        'house_number',
        'price',
        'total_area',
        'room_number',
        'floor',
        'total_floor',
    ]

    def __init__(
        self,
        url: str,
        city: str,
        street: str,
        house_number: str,
        price: float,
        total_area: float,
        room_number: int,
        floor: int,
        total_floor: int
    ):
        self.url = url
        self.city = city
        self.street = street
        self.house_number = house_number
        self.price = price
        self.total_area = total_area
        self.room_number = room_number
        self.floor = floor
        self.total_floor = total_floor

    def address(self) -> str:
        return ', '.join(
            s for s in [self.city, self.street, self.house_number] if s != ''
        )


def _read_flats(query: Query) -> DataFrame:
    with _engine.connect() as connection:
        return read_sql(
            '''
            select url,
            street,
            house_number,
            price       as actual_price,
            %s          as utmost_price,
            total_area,
            living_area,
            kitchen_area,
            room_number as actual_room_number,
            %s          as desired_room_number,
            floor       as actual_floor,
            total_floor,
            %s          as desired_floor,
            case
                when housing = 'primary' then 0
                else 1
                end     as housing,
            ssf,
            izf,
            gzf
            from flats
            where city = %s
            ''',
            connection,
            params=[query.price, query.room_number, query.floor, query.city]
        )
