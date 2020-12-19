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
        housing: int,
        longitude: float,
        latitude: float,
        city: str,
        street: str,
        house_number: str,
        ssf: float,
        izf: float,
        gzf: float,
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
        self.housing = Housing(housing)
        self.longitude = longitude
        self.latitude = latitude
        self.city = city
        self.street = street
        self.house_number = house_number
        self.ssf = ssf
        self.izf = izf
        self.gzf = gzf
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


class Image:
    __slots__ = ['url', 'label']

    def __init__(self, url: str, label: 'Label'):
        self.url = url
        self.label = label


@unique
class Label(Enum):
    abandoned = -2
    unknown = -1
    luxury = 0
    comfort = 1
    junk = 2
    construction = 3
    excess = 4
