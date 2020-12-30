from typing import Optional
from flask import Flask, request, render_template, abort
from sqlalchemy import create_engine
from rampart.config import get_config
from rampart.logging import get_handler
from rampart.ranking import Query, Ranker


def _main():
    config = get_config()
    app = Flask(
        'rampart.hemingway',
        template_folder=config.hemingway.template_path
    )
    app.logger.addHandler(get_handler())
    engine = create_engine(config.hemingway.dsn)
    ranker = Ranker(config.hemingway.ranker, engine)
    app.add_url_rule(
        '/',
        view_func=lambda: _get_index(ranker),
        methods=['GET']
    )
    try:
        app.run(
            '0.0.0.0',
            config.hemingway.port,
            load_dotenv=False,
            use_reloader=False
        )
    except KeyboardInterrupt:
        pass
    except Exception:  # noqa
        app.logger.exception('Hemingway got fatal error')
    finally:
        engine.dispose()


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
