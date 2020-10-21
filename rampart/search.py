from logging import info
from os import getenv
from typing import List
from sqlalchemy import create_engine
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine

_engine: Engine = create_engine(getenv('RAMPART_DATABASE_DSN'))
_booster = Booster(model_file='model.txt')


def search_flats(query: 'Query') -> List['Flat']:
    info(_read_flats(query))
    return []


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
        'price',
        'total_area',
        'room_number',
        'floor',
        'total_floor',
        'address',
    ]


def _read_flats(query: Query) -> DataFrame:
    with _engine.connect() as connection:
        return read_sql(
            '''
            select url,
            price       as actual_price,
            $1          as utmost_price,
            total_area,
            living_area,
            kitchen_area,
            room_number as actual_room_number,
            $2          as desired_room_number,
            floor       as actual_floor,
            total_floor,
            $3          as desired_floor,
            case
                when housing = 'primary' then 0
                else 1
                end     as housing,
            ssf,
            izf,
            gzf
            from flats
            where city = $4
            ''',
            connection,
            [query.price, query.room_number, query.floor, query.city]
        )
