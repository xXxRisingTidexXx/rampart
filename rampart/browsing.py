from typing import Optional
from flask import Flask, request, render_template, abort
from rampart.config import get_config
from rampart.ranking import Query, Ranker


def _main():
    config = get_config()
    ranker = Ranker(config.browsing.ranker)
    app = Flask('rampart.browsing', template_folder=config.browsing.template_path)
    app.add_url_rule('/', view_func=lambda: _get_index(ranker), methods=['GET'])
    app.run('0.0.0.0', config.browsing.port, load_dotenv=False, use_reloader=False)


def _get_index(ranker: Ranker) -> str:
    city = request.args.get('city')
    if not city:
        city = 'Київ'
    price = _float(request.args.get('price'))
    if price < 0:
        abort(400)
    floor = _int(request.args.get('floor'))
    if floor < 0 or floor > 2:
        abort(400)
    room_number = _int(request.args.get('room_number'))
    if room_number < 0 or room_number > 4:
        abort(400)
    limit = _int(request.args.get('limit'), 10)
    if limit < 1:
        abort(400)
    offset = _int(request.args.get('offset'))
    if offset < 0:
        abort(400)
    query = Query(city, price, floor, room_number, limit, offset)
    return render_template(
        'index.html',
        lower=query.lower + 1,
        upper=query.upper,
        city=city,
        price=price,
        floor=floor,
        room_number=room_number,
        limit=limit,
        offset=offset,
        previous=offset - 1,
        next=offset + 1,
        flats=ranker.rank_flats(query)
    )


def _float(value: Optional[str]) -> float:
    if not value:
        return 0
    try:
        return float(value)
    except ValueError:
        return -1


def _int(value: Optional[str], default: int = 0) -> int:
    if not value:
        return default
    try:
        return int(value)
    except ValueError:
        return -1


if __name__ == '__main__':
    _main()
