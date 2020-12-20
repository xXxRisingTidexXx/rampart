from argparse import ArgumentParser
from prometheus_client.exposition import start_http_server
from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from rampart.config import get_config
from rampart.logging import get_logger
from rampart.metrics import Drain
from rampart.recognition import Recognizer
from requests import Session
from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.cron import CronTrigger

_logger = get_logger('rampart.auge')


def _main():
    parser = ArgumentParser(description='Rampart image classification job.')
    parser.add_argument(
        '-debug',
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
    recognizer = Recognizer(
        config.auge.recognizer,
        engine,
        session,
        Drain(engine)
    )
    try:
        if args.debug:
            recognizer()
        else:
            start_http_server(config.auge.metrics_port)
            scheduler = BlockingScheduler()
            scheduler.add_job(
                recognizer,
                CronTrigger.from_crontab(config.auge.spec)
            )
            scheduler.start()
    except KeyboardInterrupt:
        pass
    except Exception:  # noqa
        _logger.exception('Auge got fatal error')
    finally:
        engine.dispose()
        session.close()


if __name__ == '__main__':
    _main()
