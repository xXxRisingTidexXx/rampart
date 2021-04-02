from os import getenv
from pathlib import Path
from typing import Any, Dict
from yaml import safe_load

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['auge', 'twinkle', 'coquus']

    def __init__(self, config: Dict[str, Any]):
        dsn = _get_env('RAMPART_DSN')
        config['auge']['dsn'] = dsn
        config['twinkle']['dsn'] = dsn
        self.auge = AugeConfig(config['auge'])
        self.twinkle = TwinkleConfig(config['twinkle'])
        self.coquus = CoquusConfig(config['coquus'])


def _get_env(key: str) -> str:
    value = getenv(key)
    if not value:
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value


class AugeConfig:
    __slots__ = [
        'dsn',
        'loader_number',
        'retry_limit',
        'loader',
        'model_path',
        'metrics_port',
        'buffer_size',
        'interval'
    ]

    def __init__(self, config: Dict[str, Any]):
        self.dsn: str = config['dsn']
        self.loader_number: int = config['loader-number']
        self.retry_limit: int = config['retry-limit']
        self.loader = LoaderConfig(config['loader'])
        self.model_path = str(_root_path / config['model-path'])
        self.metrics_port: int = config['metrics-port']
        self.buffer_size: int = config['buffer-size']
        self.interval: int = config['interval']


class LoaderConfig:
    __slots__ = ['timeout', 'user_agent']

    def __init__(self, config: Dict[str, Any]):
        self.timeout: float = config['timeout']
        self.user_agent: str = config['user-agent']


class TwinkleConfig:
    __slots__ = ['dsn', 'model_path', 'metrics_port', 'spec']

    def __init__(self, config: Dict[str, Any]):
        self.dsn: str = config['dsn']
        self.model_path = str(_root_path / config['model-path'])
        self.metrics_port: int = config['metrics-port']
        self.spec: str = config['spec']


class CoquusConfig:
    __slots__ = ['input_path', 'output_format']

    def __init__(self, config: Dict[str, Any]):
        self.input_path = str(_root_path / config['input-path'])
        self.output_format = str(_root_path / config['output-format'])
