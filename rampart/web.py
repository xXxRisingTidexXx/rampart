from os import getenv
from fastapi import FastAPI, HTTPException
from requests import codes
from sqlalchemy import create_engine
from rampart.twinkle import Reader, Classifier, Query, RoomNumber, Floor

app = FastAPI()
_reader = Reader(create_engine(getenv('RAMPART_DSN')))
_classifier = Classifier()


@app.get('/')
def _get_root(city: str, price: float, room_number: str, floor: str):
    if price < 0:
        raise HTTPException(codes.bad, 'Price shouldn\'t be negative')
    if room_number not in RoomNumber.__members__:
        raise HTTPException(codes.bad, 'Choose room number from the existing ones')
    if floor not in Floor.__members__:
        raise HTTPException(codes.bad, 'Choose floor from the existing ones')
    flats = _reader.read_flats(Query(city, price, RoomNumber[room_number], Floor[floor]))
    if len(flats) <= 0:
        return []
    flats['is_relevant'] = _classifier.classify_flats(
        flats.drop(columns=['id', 'url', 'street', 'house_number'])
    )
    return [
        {
            'url': f['url'],
            'price': int(f['actual_price']),
            'total_area': int(f['total_area']),
            'living_area': int(f['living_area']),
            'kitchen_area': int(f['kitchen_area']),
            'room_number': f['actual_room_number'],
            'floor': f['actual_floor'],
            'total_floor': f['total_floor'],
            'city': city,
            'street': f['street'],
            'house_number': f['house_number'],
            'ssf': f['ssf'],
            'izf': f['izf'],
            'gzf': f['gzf'],
            'abandoned_count': f['abandoned_count'],
            'luxury_count': f['luxury_count'],
            'comfort_count': f['comfort_count'],
            'junk_count': f['junk_count'],
            'construction_count': f['construction_count'],
            'excess_count': f['excess_count']
        }
        for _, f in flats[flats['is_relevant']].iterrows()
    ]
