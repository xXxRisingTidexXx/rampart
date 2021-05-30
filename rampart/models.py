from typing import List
from dataclasses import dataclass
from enum import Enum, unique
from shapely.geometry import Point


@dataclass(frozen=True)
class Flat:
    url: str
    images: List['Image']
    price: float
    total_price: float
    living_price: float
    kitchen_price: float
    room_number: int
    floor: int
    total_floor: int
    housing: 'Housing'
    point: Point
    city: str
    street: str
    house_number: str
    ssf: float
    izf: float
    gzf: float


@dataclass(frozen=True)
class Image:
    url: str
    interior: 'Interior'


@unique
class Interior(Enum):
    unknown = 0
    luxury = 1
    comfort = 2
    junk = 3
    construction = 4
    excess = 5


@unique
class Housing(Enum):
    primary = 0
    secondary = 1
