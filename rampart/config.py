from os import getenv
from pathlib import Path
from yaml import safe_load

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['browsing']

    def __init__(self, config):
        self.browsing = BrowsingConfig(config['browsing'])


class BrowsingConfig:
    __slots__ = ['port', 'template_path', 'searcher']

    def __init__(self, config):
        self.port: int = config['port']
        self.template_path = _root_path / 'templates'
        self.searcher = SearcherConfig()


class SearcherConfig:
    __slots__ = ['dsn', 'model_path']

    def __init__(self):
        self.dsn = _get_env('RAMPART_DATABASE_DSN')
        self.model_path = str(_root_path / 'scientific/models/twinkle.latest.txt')


def _get_env(key: str) -> str:
    value = getenv(key)
    if not value:
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value
