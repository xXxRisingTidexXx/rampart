from enum import Enum, unique


class Flat:
    __slots__ = [
        'url',
        'price',
        'total_area',
        'room_number',
        'floor',
        'total_floor',
        'housing',
        'city',
        'street',
        'house_number'
    ]

    def __init__(
        self,
        url: str,
        price: float,
        total_area: float,
        room_number: int,
        floor: int,
        total_floor: int,
        housing: int,
        city: str,
        street: str,
        house_number: str
    ):
        self.url = url
        self.price = price
        self.total_area = total_area
        self.room_number = room_number
        self.floor = floor
        self.total_floor = total_floor
        self.housing = Housing(housing)
        self.city = city
        self.street = street
        self.house_number = house_number

    @property
    def address(self) -> str:
        return ', '.join(
            s for s in [self.city, self.street, self.house_number] if s != ''
        )


@unique
class Housing(Enum):
    primary = 0
    secondary = 1
