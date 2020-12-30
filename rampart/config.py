from os import getenv
from pathlib import Path
from typing import Any, Dict
from yaml import safe_load

_root_path = Path(__file__).parent.parent


def get_config() -> 'Config':
    with open(_root_path / 'config/dev.yaml') as stream:
        return Config(safe_load(stream))


class Config:
    __slots__ = ['browsing', 'coquus', 'auge']

    def __init__(self, config: Dict[str, Any]):
        dsn = _get_env('RAMPART_DATABASE_DSN')
        self.browsing = BrowsingConfig(config['browsing'], dsn)
        self.coquus = CoquusConfig(config['coquus'])
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


# TODO: move model path to config.
class RankerConfig:
    __slots__ = ['dsn', 'model_path']

    def __init__(self, dsn: str):
        self.dsn = dsn
        self.model_path = str(
            _root_path / 'scientific/models/twinkle.latest.txt'
        )


class CoquusConfig:
    __slots__ = ['input_path', 'output_format']

    def __init__(self, config: Dict[str, Any]):
        self.input_path = str(_root_path / config['input-path'])
        self.output_format = str(_root_path / config['output-format'])


class AugeConfig:
    __slots__ = [
        'dsn',
        'pool_size',
        'retry_limit',
        'recognizer',
        'metrics_port',
        'spec'
    ]

    def __init__(self, config: Dict[str, Any], dsn: str):
        self.dsn = dsn
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
