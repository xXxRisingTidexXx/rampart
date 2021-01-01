from argparse import ArgumentParser
from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.cron import CronTrigger
from prometheus_client.exposition import start_http_server
from sqlalchemy import create_engine
from rampart.config import get_config
from rampart.logging import get_logger

_logger = get_logger('rampart.twinkle')


def _main():
    parser = ArgumentParser(description='Rampart flat ranking task.')
    parser.add_argument(
        '-debug',
        default=False,
        action='store_true',
        help='Whether to run the job immediately or periodically'
    )
    args = parser.parse_args()
    config = get_config()
    engine = create_engine(config.twinkle.dsn)
    try:
        if args.debug:
            _hello()
        else:
            start_http_server(config.twinkle.metrics_port)
            scheduler = BlockingScheduler()
            scheduler.add_job(
                _hello,
                CronTrigger.from_crontab(config.twinkle.spec)
            )
            scheduler.start()
    except KeyboardInterrupt:
        pass
    except Exception:  # noqa
        _logger.exception('Twinkle got fatal error')
    finally:
        engine.dispose()


def _hello():
    print('Hello world!')


if __name__ == '__main__':
    _main()