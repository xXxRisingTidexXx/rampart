from argparse import ArgumentParser
from queue import Queue
from threading import Thread
from prometheus_client.exposition import start_http_server
from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from rampart.config import get_config, AugeConfig
from rampart.logging import get_logger
from rampart.recognition import Reader, Loader, Recognizer, Updater
from requests import Session, RequestException
from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.interval import IntervalTrigger

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
    reader = Reader(engine)
    session = Session()
    session.mount(
        'https://',
        HTTPAdapter(
            pool_maxsize=config.auge.loader_number,
            max_retries=config.auge.retry_limit
        )
    )
    loader = Loader(config.auge.loader, session)
    recognizer = Recognizer(config.auge.model_path)
    updater = Updater(engine)
    try:
        if args.debug:
            _run_once(reader, loader, recognizer, updater)
        else:
            _run_forever(config.auge, reader, loader, recognizer, updater)
    except Exception:  # noqa
        _logger.exception('Auge got fatal error')
    finally:
        session.close()
        engine.dispose()


def _run_once(reader: Reader, loader: Loader, recognizer: Recognizer, updater: Updater):
    for url in reader.read_urls(1):
        updater.update_image(recognizer.recognize_image(loader.load_image(url)))


def _run_forever(
    config: AugeConfig,
    reader: Reader,
    loader: Loader,
    recognizer: Recognizer,
    updater: Updater
):
    start_http_server(config.metrics_port)
    urls = Queue(config.buffer_size)
    Thread(target=_read_urls, args=(reader, urls), daemon=True).start()
    images = Queue(config.buffer_size)
    for _ in range(config.loader_number):
        Thread(target=_load_images, args=(loader, urls, images), daemon=True).start()
    Thread(
        target=_recognize_images,
        args=(recognizer, updater, images),
        daemon=True
    ).start()
    scheduler = BlockingScheduler()
    scheduler.add_job(
        _read_urls,
        IntervalTrigger(seconds=config.interval),
        (reader, config.loader_number, urls)
    )
    try:
        scheduler.start()
    except KeyboardInterrupt:
        scheduler.shutdown()
        urls.join()
        images.join()


def _read_urls(reader: Reader, limit: int, urls: Queue):
    for url in reader.read_urls(limit):
        urls.put(url)


def _load_images(loader: Loader, urls: Queue, images: Queue):
    while True:
        url = urls.get()
        try:
            images.put(loader.load_image(url))
        except (RequestException, RuntimeError):
            _logger.exception('Auge failed to load an image', extra={'url': url})
        finally:
            urls.task_done()


def _recognize_images(recognizer: Recognizer, updater: Updater, images: Queue):
    while True:
        updater.update_image(recognizer.recognize_image(images.get()))
        images.task_done()


if __name__ == '__main__':
    _main()
