from typing import List
from sqlalchemy import create_engine
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine
from rampart.config import SearcherConfig
from rampart.model import Flat


# TODO: leverage optuna to set the hyperparameters.
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
            frame.drop(columns=['url', 'street', 'house_number']),
            num_iteration=self._booster.best_iteration
        )
        return [
            Flat(
                s['url'],
                s['actual_price'],
                s['total_area'],
                s['actual_room_number'],
                s['actual_floor'],
                s['total_floor'],
                s['housing'],
                query.city,
                s['street'],
                s['house_number'],
            )
            for _, s
            in frame.sort_values('score', ascending=False).head(7).iterrows()
        ]

    def _read_flats(self, query: 'Query') -> DataFrame:
        with self._engine.connect() as connection:
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


class Query:
    __slots__ = ['city', 'price', 'floor', 'room_number']

    def __init__(self, city: str, price: float, floor: int, room_number: int):
        self.city = city
        self.price = price
        self.floor = floor
        self.room_number = room_number
