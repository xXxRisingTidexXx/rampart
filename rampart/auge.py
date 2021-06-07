from argparse import ArgumentParser
from enum import Enum, unique
from io import BytesIO
from logging import exception
from os import getenv
from queue import Queue
from threading import Thread
from typing import List
from PIL.Image import open, new, Image as Picture
from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.interval import IntervalTrigger
from requests import Session, codes, RequestException
from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from sqlalchemy.engine.base import Engine
from torch import load, no_grad, max, unsqueeze
from torchvision.transforms import Compose, ToTensor, Resize, Normalize
from torchvision.models import Inception3
from rampart.exceptions import RampartError


class Reader:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def read_urls(self, limit: int) -> List[str]:
        with self._engine.connect() as connection:
            return [
                u[0] for u in connection.execute(
                    "select url from images where interior = 'unknown' limit %s",
                    limit
                )
            ]


class Loader:
    __slots__ = ['_session']

    def __init__(self, session: Session):
        self._session = session

    def load_image(self, url: str) -> 'Image':
        response = self._session.get(url, timeout=2)
        if response.status_code != codes.ok:
            raise RampartError(f'Loader got non-ok status {response.status_code}')
        source = open(BytesIO(response.content))
        if source.mode == 'RGBA':
            canvas = new('RGBA', source.size, 'white')
            canvas.paste(source, (0, 0), source)
            source = canvas.convert('RGB')
        return Image(url, source, Interior.unknown)


class Image:
    __slots__ = ['url', 'source', 'interior']

    def __init__(self, url: str, source: Picture, interior: 'Interior'):
        self.url = url
        self.source = source
        self.interior = interior


@unique
class Interior(Enum):
    unknown = -2
    abandoned = -1
    luxury = 0
    comfort = 1
    junk = 2
    construction = 3
    excess = 4


class Recognizer:
    __slots__ = ['_transforms', '_module']

    def __init__(self, path: str):
        self._transforms = Compose(
            [
                ToTensor(),
                Resize((299, 299)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )
        self._module = Inception3(5, init_weights=False)
        self._module.load_state_dict(load(path))
        self._module.eval()

    @no_grad()
    def recognize_image(self, image: Image) -> Image:
        return Image(
            image.url,
            image.source,
            Interior(
                max(
                    self._module(unsqueeze(self._transforms(image.source), 0)),
                    1
                ).indices.item()
            )
        )


class Updater:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def update_image(self, image: Image):
        with self._engine.connect() as connection:
            connection.execute(
                'update images set interior = %s where url = %s',
                image.interior.name,
                image.url
            )


def _main():
    parser = ArgumentParser(description='Rampart image classification job.')
    parser.add_argument(
        '-debug',
        default=False,
        action='store_true',
        help='Whether to run the job immediately or periodically'
    )
    args = parser.parse_args()
    engine = create_engine(getenv('RAMPART_DSN'))
    reader = Reader(engine)
    session = Session()
    session.mount('https://', HTTPAdapter(pool_maxsize=15, max_retries=3))
    loader = Loader(session)
    recognizer = Recognizer('scientific/models/auge.latest.pth')
    updater = Updater(engine)
    try:
        if args.debug:
            _run_once(reader, loader, recognizer, updater)
        else:
            _run_forever(reader, loader, recognizer, updater)
    except Exception:  # noqa
        exception('Auge got fatal error')
    finally:
        session.close()
        engine.dispose()


def _run_once(reader: Reader, loader: Loader, recognizer: Recognizer, updater: Updater):
    for url in reader.read_urls(1):
        updater.update_image(recognizer.recognize_image(loader.load_image(url)))


def _run_forever(
    reader: Reader,
    loader: Loader,
    recognizer: Recognizer,
    updater: Updater
):
    urls, images = Queue(300), Queue(300)
    for _ in range(15):
        Thread(target=_load_images, args=(loader, urls, images), daemon=True).start()
    Thread(
        target=_recognize_images,
        args=(recognizer, updater, images),
        daemon=True
    ).start()
    scheduler = BlockingScheduler()
    scheduler.add_job(_read_urls, IntervalTrigger(seconds=15), (reader, urls))
    try:
        scheduler.start()
    except KeyboardInterrupt:
        scheduler.shutdown()
        urls.join()
        images.join()


def _read_urls(reader: Reader, urls: Queue):
    for url in reader.read_urls(15):
        urls.put(url)


def _load_images(loader: Loader, urls: Queue, images: Queue):
    while True:
        url = urls.get()
        try:
            images.put(loader.load_image(url))
        except (RequestException, RampartError):
            exception('Auge failed to load an image', extra={'url': url})
        finally:
            urls.task_done()


def _recognize_images(recognizer: Recognizer, updater: Updater, images: Queue):
    while True:
        updater.update_image(recognizer.recognize_image(images.get()))
        images.task_done()


if __name__ == '__main__':
    _main()
