from enum import Enum, unique
from lightgbm import Booster
from numpy import ndarray
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine


class Reader:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def read_flats(self, query: 'Query') -> DataFrame:
        with self._engine.connect() as connection:
            return read_sql(
                '''
                select flats.id,
                    flats.url,
                    price        as actual_price,
                    %s           as utmost_price,
                    total_area,
                    living_area,
                    kitchen_area,
                    room_number  as actual_room_number,
                    %s           as desired_room_number,
                    floor        as actual_floor,
                    total_floor,
                    %s           as desired_floor,
                    case
                        when housing = 'primary' then 0
                        else 1
                        end      as housing,
                    ssf,
                    izf,
                    gzf,
                    sum(
                        case
                            when interior = 'abandoned' then 1
                            else 0
                            end) as abandoned_count,
                    sum(
                        case
                            when interior = 'luxury' then 1
                            else 0
                            end) as luxury_count,
                    sum(
                        case
                            when interior = 'comfort' then 1
                            else 0
                            end) as comfort_count,
                    sum(
                        case
                            when interior = 'junk' then 1
                            else 0
                            end) as junk_count,
                    sum(
                        case
                            when interior = 'construction' then 1
                            else 0
                            end) as construction_count,
                    sum(
                        case
                            when interior = 'excess' then 1
                            else 0
                            end) as excess_count,
                    0 as panorama_count,
                    street,
                    house_number
                from flats
                    join visuals on flats.id = visuals.flat_id
                    join images on visuals.image_id = images.id
                where city = %s and interior != 'unknown'
                group by flats.id
                ''',
                connection,
                params=[
                    query.price,
                    query.room_number.value,
                    query.floor.value,
                    query.city
                ]
            )


class Query:
    __slots__ = ['city', 'price', 'room_number', 'floor']

    def __init__(
        self,
        city: str,
        price: float,
        room_number: 'RoomNumber',
        floor: 'Floor'
    ):
        self.city = city
        self.price = price
        self.room_number = room_number
        self.floor = floor


@unique
class RoomNumber(Enum):
    any = 0
    one = 1
    two = 2
    three = 3
    many = 4


@unique
class Floor(Enum):
    any = 0
    low = 1
    high = 2


class Classifier:
    __slots__ = ['_booster']

    def __init__(self):
        self._booster = Booster(model_file='scientific/models/twinkle.latest.txt')

    def classify_flats(self, flats: DataFrame) -> ndarray:
        return self._booster.predict(flats).round(0).astype(bool)
