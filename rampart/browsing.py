from flask import Flask, request
from rampart.config import BrowsingHandlerConfig, get_config
from rampart.search import Query, Searcher


class Handler:
    def __init__(self, config: BrowsingHandlerConfig):
        self._searcher = Searcher(config.searcher)

    def get_root(self):
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
        return 'Hello, world!'


if __name__ == '__main__':
    setup = get_config().browsing
    app = Flask('rampart.browsing')
    handler = Handler(setup.handler)
    app.add_url_rule('/', view_func=handler.get_root, methods=['GET'])
    app.run('0.0.0.0', setup.port, load_dotenv=False, use_reloader=False)
