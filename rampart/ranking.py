from typing import List
from enum import Enum, unique
from lightgbm import Booster
from pandas import read_sql, DataFrame
from sqlalchemy.engine.base import Engine
from rampart.config import RankerConfig


# TODO: leverage optuna to set the hyperparameters.
# https://medium.com/optuna/lightgbm-tuner-new-optuna-integration-for-hyperparameter-optimization-8b7095e99258
# https://scikit-learn.org/stable/modules/generated/sklearn.metrics.ndcg_score.html
# https://lightgbm.readthedocs.io/en/latest/Parameters-Tuning.html
# TODO: add ranking metrics.
class Ranker:
    __slots__ = ['_loader', '_reader', '_booster', '_writer', '_limit']

    def __init__(self, config: RankerConfig, engine: Engine):
        self._loader = Loader(engine)
        self._reader = Reader(engine, config.price_factor)
        self._booster = Booster(model_file=config.model_path)
        self._writer = Writer(engine)
        self._limit = config.limit

    def __call__(self):
        for subscription in self._loader.load_subscriptions():
            flats = self._reader.read_flats(subscription)
            if len(flats) > 0:
                flats['score'] = self._booster.predict(
                    flats.drop(columns=['id']),
                    num_iteration=self._booster.best_iteration
                )
                self._writer.write_lookup(
                    Lookup(
                        subscription.id,
                        flats
                        .sort_values('score', ascending=False)
                        .head(self._limit)['id']
                        .tolist()
                    )
                )


class Loader:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def load_subscriptions(self) -> List['Subscription']:
        with self._engine.connect() as connection:
            return [
                Subscription(s[0], s[1], s[2], RoomNumber[s[3]], Floor[s[4]])
                for s in connection.execute(
                    '''
                    select id, city, price, room_number, floor
                    from subscriptions
                    '''
                )
            ]


class Subscription:
    __slots__ = ['id', 'city', 'price', 'room_number', 'floor']

    def __init__(
        self,
        id_: int,
        city: str,
        price: float,
        room_number: 'RoomNumber',
        floor: 'Floor'
    ):
        self.id = id_
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


class Reader:
    __slots__ = ['_engine', '_price_factor']

    def __init__(self, engine: Engine, price_factor: float):
        self._engine = engine
        self._price_factor = price_factor

    def read_flats(self, subscription: Subscription) -> DataFrame:
        price_clause = ''
        if subscription.price > 0:
            price_clause = (
                f'and price <= {self._price_factor * subscription.price}'
            )
        room_number_clause = ''
        if subscription.room_number == RoomNumber.many:
            room_number_clause = f'and room_number >= {RoomNumber.many.value}'
        elif subscription.room_number != RoomNumber.any:
            room_number_clause = (
                f'and room_number = {subscription.room_number.value}'
            )
        with self._engine.connect() as connection:
            return read_sql(
                f'''
                select flats.id,
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
                    and flats.id not in (
                        select entries.flat_id
                        from entries
                            join lookups on entries.lookup_id = lookups.id
                        where subscription_id = %s)
                {price_clause}
                {room_number_clause}
                group by flats.id
                having sum(
                    case
                        when kind = 'photo' and label = 'unknown' then 1
                        else 0
                        end) = 0
                ''',
                connection,
                params=[
                    subscription.price,
                    subscription.room_number.value,
                    subscription.floor.value,
                    subscription.city,
                    subscription.id
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
                insert into lookups (subscription_id, creation_time, status)
                values (%s, now() at time zone 'utc', 'unseen')
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
