from os import getenv
from re import compile
from pathlib import Path
from typing import Dict
from yaml import safe_load
from rampart.model import Housing

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['telebot']

    def __init__(self, config):
        self.telebot = TelebotConfig(config['telebot'])


class TelebotConfig:
    __slots__ = ['token', 'start_command', 'help_command', 'search_pattern', 'server']

    def __init__(self, config):
        self.token = _get_env('RAMPART_TELEBOT_TOKEN')
        self.start_command: str = config['start_command']
        self.help_command: str = config['help_command']
        self.search_pattern = compile(config['search_pattern'])
        self.server = ServerConfig(config['server'])


def _get_env(key: str) -> str:
    value = getenv(key)
    if value != '':
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value


class ServerConfig:
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
        self.room_numbers: Dict[str, int] = config['room_numbers']
        self.housings: Dict[Housing, str] = {
            Housing(i): h for i, h in enumerate(config['housings'])
        }
        self.start_template = _read_template('start')
        self.help_template = _read_template('help')
        self.nothing_template = _read_template('nothing')
        self.flat_template = _read_template('flat')
        self.confusion_template = _read_template('confusion')
        self.searcher = SearcherConfig()


def _read_template(name: str) -> str:
    with open(_root_path / f'templates/{name}.html') as stream:
        return stream.read()


class SearcherConfig:
    __slots__ = ['dsn', 'model_path']

    def __init__(self):
        self.dsn = _get_env('RAMPART_DATABASE_DSN')
        self.model_path = _root_path / 'model.txt'
