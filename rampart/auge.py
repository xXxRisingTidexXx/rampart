from argparse import ArgumentParser
from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from rampart.config import get_config
from rampart.logging import get_logger
from rampart.recognition import Recognizer
from requests import Session
from schedule import every, run_pending
from time import sleep

_logger = get_logger('rampart.auge')


def _main():
    parser = ArgumentParser(description='Rampart image classification job.')
    parser.add_argument(
        '-dev',
        default=False,
        action='store_true',
        help='Whether to run the job immediately or periodically'
    )
    args = parser.parse_args()
    config = get_config()
    engine = create_engine(config.auge.dsn)
    session = Session()
    session.mount(
        'https://',
        HTTPAdapter(
            pool_maxsize=config.auge.pool_size,
            max_retries=config.auge.retry_limit
        )
    )
    recognizer = Recognizer(config.auge.recognizer, engine, session)
    try:
        if args.dev:
            recognizer()
        else:
            every(config.auge.interval).minutes.do(recognizer)
            while True:
                run_pending()
                sleep(1)
    except KeyboardInterrupt:
        pass
    except Exception:  # noqa
        _logger.exception('Auge got fatal error')
    finally:
        engine.dispose()
        session.close()


if __name__ == '__main__':
    _main()
