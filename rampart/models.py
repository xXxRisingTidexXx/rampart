from typing import List
from dataclasses import dataclass
from enum import Enum, unique
from shapely.geometry import Point


@dataclass(frozen=True)
class Flat:
    url: str
    images_urls: List[str]
    price: float
    total_area: float
    living_area: float
    kitchen_area: float
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

    @property
    def has_location(self) -> bool:
        return self.point.x or self.point.y


@unique
class Housing(Enum):
    primary = 0
    secondary = 1


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
