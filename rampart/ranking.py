from typing import List
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine
from rampart.config import RankerConfig
from rampart.models import Flat


# TODO: leverage optuna to set the hyperparameters.
class Ranker:
    __slots__ = ['_reader', '_booster']

    def __init__(self, config: RankerConfig, engine: Engine):
        self._reader = Reader(engine)
        self._booster = Booster(model_file=config.model_path)

    def rank_flats(self, query: 'Query') -> List[Flat]:
        frame = self._reader.read_flats(query)
        if len(frame) == 0:
            return []
        frame['score'] = self._booster.predict(
            frame.drop(
                columns=[
                    'id',
                    'url',
                    'longitude',
                    'latitude',
                    'street',
                    'house_number'
                ]
            ),
            num_iteration=self._booster.best_iteration
        )
        return [
            Flat(
                s['id'],
                s['url'],
                s['actual_price'],
                s['total_area'],
                s['living_area'],
                s['kitchen_area'],
                s['actual_room_number'],
                s['actual_floor'],
                s['total_floor'],
                s['housing'],
                s['longitude'],
                s['latitude'],
                query.city,
                s['street'],
                s['house_number'],
                s['ssf'],
                s['izf'],
                s['gzf'],
                s['score']
            )
            for _, s
            in (
                frame
                .sort_values('score', ascending=False)
                .iloc[query.lower:query.upper]
                .iterrows()
            )
        ]


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
                       st_x(point)  as longitude,
                       st_y(point)  as latitude,
                       street,
                       house_number,
                       ssf,
                       izf,
                       gzf,
                       sum(
                           case
                               when kind = 'photo' and label = 'unknown' then 1
                               else 0
                               end) as unknown_count,
                       sum(
                           case
                               when kind = 'photo'
                                    and label = 'abandoned' then 1
                               else 0
                               end) as abandoned_count,
                       sum(
                           case
                               when kind = 'photo' and label = 'luxury' then 1
                               else 0
                               end) as luxury_count,
                       sum(
                           case
                               when kind = 'photo' and label = 'comfort' then 1
                               else 0
                               end) as comfort_count,
                       sum(
                           case
                               when kind = 'photo' and label = 'junk' then 1
                               else 0
                               end) as junk_count,
                       sum(
                           case
                               when kind = 'photo'
                                    and label = 'construction' then 1
                               else 0
                               end) as construction_count,
                       sum(
                           case
                               when kind = 'photo' and label = 'excess' then 1
                               else 0
                               end) as excess_count,
                       sum(
                           case
                               when kind = 'panorama' then 1
                               else 0
                               end) as panorama_count
                from flats
                     join images on flats.id = flat_id
                where city = %s
                group by flats.id
                ''',
                connection,
                params=[
                    query.price,
                    query.room_number,
                    query.floor,
                    query.city
                ]
            )


class Query:
    __slots__ = ['city', 'price', 'floor', 'room_number', 'limit', 'offset']

    def __init__(
        self,
        city: str,
        price: float,
        floor: int,
        room_number: int,
        limit: int,
        offset: int
    ):
        self.city = city
        self.price = price
        self.floor = floor
        self.room_number = room_number
        self.limit = limit
        self.offset = offset

    @property
    def lower(self) -> int:
        return self.limit * self.offset

    @property
    def upper(self) -> int:
        return self.limit * (self.offset + 1)
