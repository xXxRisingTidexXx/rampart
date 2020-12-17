from os import getenv
from pathlib import Path
from typing import Any, Dict
from yaml import safe_load

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['browsing', 'auge']

    def __init__(self, config: Dict[str, Any]):
        dsn = _get_env('RAMPART_DATABASE_DSN')
        self.browsing = BrowsingConfig(config['browsing'], dsn)
        self.auge = AugeConfig(config['auge'], dsn)


def _get_env(key: str) -> str:
    value = getenv(key)
    if not value:
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value


class BrowsingConfig:
    __slots__ = ['port', 'template_path', 'ranker']

    def __init__(self, config: Dict[str, Any], dsn: str):
        self.port: int = config['port']
        self.template_path = _root_path / 'templates'
        self.ranker = RankerConfig(dsn)


class RankerConfig:
    __slots__ = ['dsn', 'model_path']

    def __init__(self, dsn: str):
        self.dsn = dsn
        self.model_path = str(_root_path / 'scientific/models/twinkle.latest.txt')


class AugeConfig:
    __slots__ = ['thread_number', 'retry_limit', 'dsn', 'model_path']

    def __init__(self, config: Dict[str, Any], dsn: str):
        self.thread_number: int = config['thread-number']
        self.retry_limit: int = config['retry-limit']
        self.dsn = dsn
        self.model_path = str(_root_path / 'scientific/models/auge.latest.pth')
