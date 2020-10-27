from typing import Tuple, Union, Optional
from flask import Flask, request, render_template
from jinja2 import Template
from rampart.config import BrowsingHandlerConfig, get_config
from rampart.search import Query, Searcher


class Handler:
    __slots__ = ['_index_name', '_searcher']

    def __init__(self, config: BrowsingHandlerConfig):
        self._index_name = config.index_name
        self._searcher = Searcher(config.searcher)

    def get_index(self) -> Tuple[Union[str, Template], int]:
        city = request.args.get('city')
        if not city:
            city = 'Київ'
        price = _float(request.args.get('price'))
        if price < 0:
            return 'Invalid price provided', 400
        floor = _int(request.args.get('floor'))
        if floor < 0 or floor > 2:
            return 'Invalid floor provided', 400
        room_number = _int(request.args.get('room_number'))
        if room_number < 0 or room_number > 4:
            return 'Invalid floor provided', 400
        limit = _int(request.args.get('limit'))
        if limit < 1:
            return 'Invalid limit provided', 400
        offset = _int(request.args.get('offset'))
        if offset < 0:
            return 'Invalid offset provided', 400
        print(
            len(
                self._searcher.search_flats(
                    Query(city, price, floor, room_number, limit, offset)
                )
            )
        )
        return render_template(self._index_name), 200


def _float(value: Optional[str]) -> float:
    if not value:
        return 0
    try:
        return float(value)
    except ValueError:
        return -1


def _int(value: Optional[str]) -> int:
    if not value:
        return 0
    try:
        return int(value)
    except ValueError:
        return -1


if __name__ == '__main__':
    setup = get_config().browsing
    app = Flask('rampart.browsing', template_folder=setup.template_path)
    handler = Handler(setup.handler)
    app.add_url_rule('/', view_func=handler.get_index, methods=['GET'])
    app.run('0.0.0.0', setup.port, load_dotenv=False, use_reloader=False)
