from os import getenv
from fastapi import FastAPI, Request
from sqlalchemy import create_engine
from fastapi.templating import Jinja2Templates
from fastapi.responses import HTMLResponse
from rampart.classification import Reader, Classifier, Query, RoomNumber, Floor

app = FastAPI()
_templates = Jinja2Templates('templates')
_reader = Reader(create_engine(getenv('RAMPART_DSN')))
_classifier = Classifier()


@app.get('/', response_class=HTMLResponse)
def _read_root(
    request: Request,
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
        return _templates.TemplateResponse('flats.html', {'request': request, 'flats': []})
    flats['is_relevant'] = _classifier.classify_flats(
        flats.drop(columns=['id', 'url', 'street', 'house_number'])
    )
    return _templates.TemplateResponse(
        'flats.html',
        {
            'request': request,
            'flats': [
                {
                    'url': f['url'],
                    'price': int(f['actual_price']),
                    'room_number': f['actual_room_number'],
                    'floor': f['actual_floor'],
                    'total_floor': f['total_floor'],
                    'city': city,
                    'street': f['street'],
                    'house_number': f['house_number']
                }
                for _, f in flats[flats['is_relevant']].iterrows()
            ]
        }
    )
