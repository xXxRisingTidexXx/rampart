from typing import List
from enum import Enum, unique
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine
from rampart.config import RankerConfig
from rampart.models import Flat, Housing


# TODO: leverage optuna to set the hyperparameters.
# https://medium.com/optuna/lightgbm-tuner-new-optuna-integration-for-hyperparameter-optimization-8b7095e99258
# https://scikit-learn.org/stable/modules/generated/sklearn.metrics.ndcg_score.html
# https://lightgbm.readthedocs.io/en/latest/Parameters-Tuning.html
# TODO: ignore flats with unknown images to avoid processing of unchecked publications.
# TODO: remove unknown count feature.
# TODO: move selection length to config.
class Ranker:
    __slots__ = ['_loader', '_reader', '_booster', '_writer']

    def __init__(self, config: RankerConfig, engine: Engine):
        self._loader = Loader(engine)
        self._reader = Reader(engine, config.price_factor)
        self._booster = Booster(model_file=config.model_path)
        self._writer = Writer(engine)

    def __call__(self):
        for subscription in self._loader.load_subscriptions():
            self._writer.write_lookup(
                Lookup(
                    subscription.id,
                    [f.id for f in self.rank_flats(subscription.query)]
                )
            )

    def rank_flats(self, query: 'Query') -> List[Flat]:
        flats = self._reader.read_flats(query)
        if len(flats) == 0:
            return []
        flats['score'] = self._booster.predict(
            flats.drop(
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
                Housing(s['housing']),
                s['longitude'],
                s['latitude'],
                query.city,
                s['street'],
                s['house_number'],
                s['ssf'],
                s['izf'],
                s['gzf'],
                s['unknown_count'],
                s['abandoned_count'],
                s['luxury_count'],
                s['comfort_count'],
                s['junk_count'],
                s['construction_count'],
                s['excess_count'],
                s['panorama_count'],
                s['score']
            )
            for _, s
            in (
                flats
                .sort_values('score', ascending=False)
                .iloc[query.lower:query.upper]
                .iterrows()
            )
        ]


class Loader:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def load_subscriptions(self) -> List['Subscription']:
        with self._engine.connect() as connection:
            return [
                Subscription(
                    s[0],
                    Query(s[1], s[2], RoomNumber[s[3]], Floor[s[4]], 3, 0)
                )
                for s in connection.execute(
                    '''
                    select id, city, price, room_number, floor
                    from subscriptions
                    '''
                )
            ]


class Subscription:
    __slots__ = ['id', 'query']

    def __init__(self, id_: int, query: 'Query'):
        self.id = id_
        self.query = query


class Query:
    __slots__ = ['city', 'price', 'room_number', 'floor', 'limit', 'offset']

    def __init__(
        self,
        city: str,
        price: float,
        room_number: 'RoomNumber',
        floor: 'Floor',
        limit: int,
        offset: int
    ):
        self.city = city
        self.price = price
        self.room_number = room_number
        self.floor = floor
        self.limit = limit
        self.offset = offset

    @property
    def lower(self) -> int:
        return self.limit * self.offset

    @property
    def upper(self) -> int:
        return self.limit * (self.offset + 1)


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


class Reader:
    __slots__ = ['_engine', '_price_factor']

    def __init__(self, engine: Engine, price_factor: float):
        self._engine = engine
        self._price_factor = price_factor

    def read_flats(self, query: 'Query') -> DataFrame:
        price_clause = ''
        if query.price > 0:
            price_clause = f'and price <= {self._price_factor * query.price}'
        room_number_clause = ''
        if query.room_number == RoomNumber.many:
            room_number_clause = f'and room_number >= {RoomNumber.many.value}'
        elif query.room_number != RoomNumber.any:
            room_number_clause = f'and room_number = {query.room_number.value}'
        with self._engine.connect() as connection:
            return read_sql(
                f'''
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
                {price_clause}
                {room_number_clause}
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


class Writer:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def write_lookup(self, lookup: 'Lookup'):
        with self._engine.begin() as connection:
            id_ = connection.scalar(
                '''
                insert into lookups (subscription_id, status)
                values (%s, 'unseen')
                returning id
                ''',
                lookup.subscription_id
            )
            connection.execute(
                '''
                insert into entries (lookup_id, flat_id, position)
                values (%s, %s, %s)
                ''',
                *[(id_, f, i) for i, f in enumerate(lookup.flat_ids)]
            )


class Lookup:
    __slots__ = ['subscription_id', 'flat_ids']

    def __init__(self, subscription_id: int, flat_ids: List[int]):
        self.subscription_id = subscription_id
        self.flat_ids = flat_ids
