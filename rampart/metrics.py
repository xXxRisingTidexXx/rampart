from time import time
from typing import Dict
from sqlalchemy.engine.base import Engine
from enum import Enum, unique


class Drain:
    __slots__ = ['_engine', '_numbers', '_durations']

    def __init__(self, engine: Engine):
        self._engine = engine
        self._numbers: Dict[Number, int] = {n: 0 for n in Number}
        self._durations: Dict[Duration, _Bucket] = {
            d: _Bucket() for d in Duration
        }

    def drain_number(self, number: 'Number'):
        self._numbers[number] += 1

    def drain_duration(self, duration: 'Duration', start: float):
        self._durations[duration].span(start)

    def flush(self):
        with self._engine.connect() as connection:
            connection.execute(
                '''
                insert into recognitions
                (
                    completion_time, abandoned_number, luxury_number,
                    comfort_number, junk_number, construction_number,
                    excess_number, reading_duration, loading_duration,
                    update_duration, total_duration
                ) values (
                    now() at time zone 'utc', %s, %s, %s, %s, %s, %s, %s, %s,
                    %s, %s 
                )
                ''',
                self._numbers[Number.abandoned],
                self._numbers[Number.luxury],
                self._numbers[Number.comfort],
                self._numbers[Number.junk],
                self._numbers[Number.construction],
                self._numbers[Number.excess],
                self._durations[Duration.reading].avg(),
                self._durations[Duration.loading].avg(),
                self._durations[Duration.update].avg(),
                self._durations[Duration.total].avg()
            )
        for number in self._numbers:
            self._numbers[number] = 0
        for bucket in self._durations.values():
            bucket.reset()


@unique
class Number(Enum):
    abandoned = -1
    luxury = 0
    comfort = 1
    junk = 2
    construction = 3
    excess = 4


@unique
class Duration(Enum):
    reading = 0
    loading = 1
    update = 2
    total = 3


class _Bucket:
    __slots__ = ['_sum', '_count']

    def __init__(self):
        self._sum = 0.0
        self._count = 0.0

    def span(self, start: float):
        self._sum += time() - start
        self._count += 1

    def avg(self) -> float:
        return 0 if self._count == 0 else self._sum / self._count

    def reset(self):
        self._sum = 0.0
        self._count = 0.0
