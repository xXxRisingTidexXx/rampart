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
    frame = _read_flats(query)
    frame['score'] = _booster.predict(
        frame.drop(columns=['url']),
        num_iteration=_booster.best_iteration
    )
    info(frame.sort_values('score', ascending=False).head())
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
    ]

    def __init__(self):
        pass


def _read_flats(query: Query) -> DataFrame:
    with _engine.connect() as connection:
        return read_sql(
            '''
            select url,
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
