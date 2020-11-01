from os import getenv
from re import compile
from pathlib import Path
from typing import Dict
from yaml import safe_load
from rampart.models import Housing

_root_path = Path(__file__).parent.parent
_template_path = _root_path / 'templates'


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['telebot', 'browsing']

    def __init__(self, config):
        self.telebot = TelebotConfig(config['telebot'])
        self.browsing = BrowsingConfig(config['browsing'])


class TelebotConfig:
    __slots__ = ['token', 'pattern', 'handler']

    def __init__(self, config):
        self.token = _get_env('RAMPART_TELEBOT_TOKEN')
        self.pattern = compile(config['pattern'])
        self.handler = TelebotHandlerConfig(config['handler'])


def _get_env(key: str) -> str:
    value = getenv(key)
    if not value:
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value


class TelebotHandlerConfig:
    __slots__ = [
        'any',
        'cities',
        'floors',
        'room_numbers',
        'housings',
        'start_template',
        'help_template',
        'nothing_template',
        'flat_template',
        'confusion_template',
        'searcher'
    ]

    def __init__(self, config):
        self.any: str = config['any']
        self.cities: Dict[str, str] = config['cities']
        self.floors: Dict[str, int] = config['floors']
        self.room_numbers: Dict[str, int] = config['room-numbers']
        self.housings: Dict[Housing, str] = {
            Housing(i): h for i, h in enumerate(config['housings'])
        }
        self.start_template = _read_template('start.html')
        self.help_template = _read_template('help.html')
        self.nothing_template = _read_template('nothing.html')
        self.flat_template = _read_template('flat.html')
        self.confusion_template = _read_template('confusion.html')
        self.searcher = SearcherConfig()


def _read_template(name: str) -> str:
    with open(_template_path / name) as stream:
        return stream.read()


class SearcherConfig:
    __slots__ = ['dsn', 'model_path']

    def __init__(self):
        self.dsn = _get_env('RAMPART_DATABASE_DSN')
        self.model_path = str(_root_path / 'model.txt')


class BrowsingConfig:
    __slots__ = ['port', 'template_path', 'handler']

    def __init__(self, config):
        self.port: int = config['port']
        self.template_path = _template_path
        self.handler = BrowsingHandlerConfig()


class BrowsingHandlerConfig:
    __slots__ = ['index_name', 'searcher']

    def __init__(self):
        self.index_name = 'index.html'
        self.searcher = SearcherConfig()
