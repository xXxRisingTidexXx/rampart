from typing import List
from sqlalchemy import create_engine
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine
from rampart.config import SearcherConfig
from rampart.models import Flat


# TODO: leverage optuna to set the hyperparameters.
# TODO: rename the module to ranking and the class to Ranker.
class Searcher:
    __slots__ = ['_engine', '_booster']

    def __init__(self, config: SearcherConfig):
        self._engine: Engine = create_engine(config.dsn)
        self._booster = Booster(model_file=config.model_path)

    def search_flats(self, query: 'Query') -> List[Flat]:
        frame = self._read_flats(query)
        if len(frame) == 0:
            return []
        frame['score'] = self._booster.predict(
            frame[
                [
                    'actual_price',
                    'utmost_price',
                    'total_area',
                    'living_area',
                    'kitchen_area',
                    'actual_room_number',
                    'desired_room_number',
                    'actual_floor',
                    'total_floor',
                    'desired_floor',
                    'housing',
                    'ssf',
                    'izf',
                    'gzf'
                ]
            ],
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

    def _read_flats(self, query: 'Query') -> DataFrame:
        with self._engine.connect() as connection:
            return read_sql(
                '''
                select id,
                       url,
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
                       st_x(point) as longitude,
                       st_y(point) as latitude,
                       street,
                       house_number,
                       ssf,
                       izf,
                       gzf
                from flats
                where city = %s
                ''',
                connection,
                params=[query.price, query.room_number, query.floor, query.city]
            )


class Query:
    __slots__ = ['city', 'price', 'floor', 'room_number', 'limit', 'offset']

    def __init__(
        self,
        city: str,
        price: float,
        floor: int,
        room_number: int,
        limit: int = 10,
        offset: int = 0
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
