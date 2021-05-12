from os import getenv
from fastapi import FastAPI
from sqlalchemy import create_engine
from rampart.classification import Reader, Classifier, Query, RoomNumber, Floor

app = FastAPI()
_reader = Reader(create_engine(getenv('RAMPART_DSN')))
_classifier = Classifier()


@app.get('/')
def _read_root(
    city: str = 'Київ',
    price: float = 0,
    room_number: str = 'any',
    floor: str = 'any',
):
    if room_number not in RoomNumber.__members__:
        room_number = 'any'
    if floor not in Floor.__members__:
        floor = 'any'
    flats = _reader.read_flats(Query(city, price, RoomNumber[room_number], Floor[floor]))
    if len(flats) <= 0:
        return []
    flats['is_relevant'] = _classifier.classify_flats(flats.drop(columns=['id', 'url']))
    return [
        {
            'url': f['url'],
            'price': f['actual_price'],
            'room_number': f['actual_room_number'],
            'floor': f['actual_floor']
        }
        for _, f in flats[flats['is_relevant']].iterrows()
    ]
