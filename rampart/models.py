from enum import Enum, unique


class Flat:
    __slots__ = [
        'id',
        'url',
        'price',
        'total_area',
        'living_area',
        'kitchen_area',
        'room_number',
        'floor',
        'total_floor',
        'housing',
        'longitude',
        'latitude',
        'city',
        'street',
        'house_number',
        'ssf',
        'izf',
        'gzf',
        'unknown_count',
        'abandoned_count',
        'luxury_count',
        'comfort_count',
        'junk_count',
        'construction_count',
        'excess_count',
        'panorama_count',
        'score'
    ]

    def __init__(
        self,
        id_: int,
        url: str,
        price: float,
        total_area: float,
        living_area: float,
        kitchen_area: float,
        room_number: int,
        floor: int,
        total_floor: int,
        housing: 'Housing',
        longitude: float,
        latitude: float,
        city: str,
        street: str,
        house_number: str,
        ssf: float,
        izf: float,
        gzf: float,
        unknown_count: int,
        abandoned_count: int,
        luxury_count: int,
        comfort_count: int,
        junk_count: int,
        construction_count: int,
        excess_count: int,
        panorama_count: int,
        score: float
    ):
        self.id = id_
        self.url = url
        self.price = price
        self.total_area = total_area
        self.living_area = living_area
        self.kitchen_area = kitchen_area
        self.room_number = room_number
        self.floor = floor
        self.total_floor = total_floor
        self.housing = housing
        self.longitude = longitude
        self.latitude = latitude
        self.city = city
        self.street = street
        self.house_number = house_number
        self.ssf = ssf
        self.izf = izf
        self.gzf = gzf
        self.unknown_count = unknown_count
        self.abandoned_count = abandoned_count
        self.luxury_count = luxury_count
        self.comfort_count = comfort_count
        self.junk_count = junk_count
        self.construction_count = construction_count
        self.excess_count = excess_count
        self.panorama_count = panorama_count
        self.score = score

    @property
    def address(self) -> str:
        return ', '.join(
            s for s in [self.city, self.street, self.house_number] if s != ''
        )


@unique
class Housing(Enum):
    primary = 0
    secondary = 1


# TODO: move to recognition.
class Image:
    __slots__ = ['url', 'label']

    def __init__(self, url: str, label: 'Label'):
        self.url = url
        self.label = label


@unique
class Label(Enum):
    unknown = -2
    abandoned = -1
    luxury = 0
    comfort = 1
    junk = 2
    construction = 3
    excess = 4
