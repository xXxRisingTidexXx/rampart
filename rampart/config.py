from os import getenv
from pathlib import Path
from typing import Any, Dict
from yaml import safe_load

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['twinkle', 'coquus', 'hemingway', 'auge']

    def __init__(self, config: Dict[str, Any]):
        dsn = _get_env('RAMPART_DATABASE_DSN')
        config['twinkle']['dsn'] = dsn
        config['twinkle']['ranker'] = RankerConfig(config['ranker'])
        config['hemingway']['dsn'] = dsn
        config['hemingway']['ranker'] = config['twinkle']['ranker']
        config['auge']['dsn'] = dsn
        self.twinkle = TwinkleConfig(config['twinkle'])
        self.coquus = CoquusConfig(config['coquus'])
        self.hemingway = HemingwayConfig(config['hemingway'])
        self.auge = AugeConfig(config['auge'])


def _get_env(key: str) -> str:
    value = getenv(key)
    if not value:
        raise RuntimeError(f'Environment variable \'{key}\' not set')
    return value


class RankerConfig:
    __slots__ = ['model_path', 'price_factor']

    def __init__(self, config: Dict[str, Any]):
        self.model_path = str(_root_path / config['model-path'])
        self.price_factor: float = config['price-factor']


class TwinkleConfig:
    __slots__ = ['dsn', 'ranker', 'metrics_port', 'spec']

    def __init__(self, config: Dict[str, Any]):
        self.dsn: str = config['dsn']
        self.ranker: RankerConfig = config['ranker']
        self.metrics_port: int = config['metrics-port']
        self.spec: str = config['spec']


class CoquusConfig:
    __slots__ = ['input_path', 'output_format']

    def __init__(self, config: Dict[str, Any]):
        self.input_path = str(_root_path / config['input-path'])
        self.output_format = str(_root_path / config['output-format'])


class HemingwayConfig:
    __slots__ = ['dsn', 'port', 'template_path', 'ranker']

    def __init__(self, config: Dict[str, Any]):
        self.dsn: str = config['dsn']
        self.port: int = config['port']
        self.template_path = str(_root_path / config['template-path'])
        self.ranker: RankerConfig = config['ranker']


class AugeConfig:
    __slots__ = [
        'dsn',
        'pool_size',
        'retry_limit',
        'recognizer',
        'metrics_port',
        'spec'
    ]

    def __init__(self, config: Dict[str, Any]):
        self.dsn: str = config['dsn']
        self.pool_size: int = config['pool-size']
        self.retry_limit: int = config['retry-limit']
        self.recognizer = RecognizerConfig(config['recognizer'])
        self.metrics_port: int = config['metrics-port']
        self.spec: str = config['spec']


class RecognizerConfig:
    __slots__ = ['model_path', 'timeout', 'batch_size', 'worker_number']

    def __init__(self, config: Dict[str, Any]):
        self.model_path = str(_root_path / config['model-path'])
        self.timeout: float = config['timeout']
        self.batch_size: int = config['batch-size']
        self.worker_number: int = config['worker-number']
