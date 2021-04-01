from argparse import ArgumentParser
from queue import Queue
from prometheus_client.exposition import start_http_server
from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from rampart.config import get_config
from rampart.logging import get_logger
from rampart.recognition import Reader, Image, Loader, Recognizer, Updater
from requests import Session
from apscheduler.schedulers.blocking import BlockingScheduler

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
    reader = Reader(engine, config.auge.loader_number)
    session = Session()
    session.mount(
        'https://',
        HTTPAdapter(
            pool_maxsize=config.auge.loader_number,
            max_retries=config.auge.retry_limit
        )
    )
    loader = Loader(session, config.auge.loader)
    recognizer = Recognizer(config.auge.model_path)
    updater = Updater(engine)
    scheduler = BlockingScheduler()
    try:
        if args.debug:
            pass
        else:
            start_http_server(config.auge.metrics_port)
            scheduler.start()
    except KeyboardInterrupt:
        scheduler.shutdown()
    except Exception:  # noqa
        _logger.exception('Auge got fatal error')
    finally:
        session.close()
        engine.dispose()


def _read_urls(reader: Reader, urls: Queue[str]):
    for url in reader.read_urls():
        urls.put(url)


def _load_images(loader: Loader, urls: Queue[str], images: Queue[Image]):
    while True:
        images.put(loader.load_image(urls.get()))
        urls.task_done()


def _recognize_images(recognizer: Recognizer, updater: Updater, images: Queue[Image]):
    while True:
        updater.update_image(recognizer.recognize_image(images.get()))
        images.task_done()


if __name__ == '__main__':
    _main()
