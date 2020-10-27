from flask import Flask, request, render_template
from rampart.config import BrowsingHandlerConfig, get_config
from rampart.search import Query, Searcher


class Handler:
    __slots__ = ['_index_name', '_searcher']

    def __init__(self, config: BrowsingHandlerConfig):
        self._index_name = config.index_name
        self._searcher = Searcher(config.searcher)

    def get_index(self):
        print(
            Query(
                request.args.get('city'),
                request.args.get('price'),
                request.args.get('floor'),
                request.args.get('room_number'),
                request.args.get('limit'),
                request.args.get('offset')
            )
        )
        return render_template(self._index_name)


if __name__ == '__main__':
    setup = get_config().browsing
    app = Flask('rampart.browsing', template_folder=setup.template_path)
    handler = Handler(setup.handler)
    app.add_url_rule('/', view_func=handler.get_index, methods=['GET'])
    app.run('0.0.0.0', setup.port, load_dotenv=False, use_reloader=False)
