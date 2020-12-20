from io import BytesIO
from time import time
from typing import List, Tuple
from PIL.Image import open, new
from requests import Session, codes
from requests.exceptions import RequestException
from sqlalchemy.engine.base import Engine
from torch.nn import (
    Module, Sequential, ReLU, Conv2d, MaxPool2d, Dropout, Linear
)
from torch import Tensor, empty, load, no_grad, max
from torch.utils.data.dataloader import default_collate, DataLoader
from torch.utils.data.dataset import Dataset
from torchvision.transforms import Compose, ToTensor, Resize, Normalize
from rampart.config import RecognizerConfig
from rampart.logging import get_logger
from rampart.metrics import Drain, Duration, Number
from rampart.models import Image, Label

_logger = get_logger('rampart.recognition')


class Recognizer:
    __slots__ = [
        '_reader',
        '_network',
        '_updater',
        '_session',
        '_timeout',
        '_batch_size',
        '_worker_number',
        '_drain'
    ]

    def __init__(
        self,
        config: RecognizerConfig,
        engine: Engine,
        session: Session,
        drain: Drain
    ):
        self._reader = Reader(engine, drain)
        self._network = Network()
        self._network.load_state_dict(load(config.model_path))
        self._network.eval()
        self._updater = Updater(engine, drain)
        self._session = session
        self._timeout = config.timeout
        self._batch_size = config.batch_size
        self._worker_number = config.worker_number
        self._drain = drain

    @no_grad()
    def __call__(self):
        start = time()
        loader = DataLoader(
            Gallery(
                self._reader.read_urls(),
                self._session,
                self._timeout,
                self._drain
            ),
            self._batch_size,
            num_workers=self._worker_number,
            collate_fn=_collate
        )
        for batch in loader:
            for url in batch[0]:
                self._updater.update_image(Image(url, Label.abandoned))
            if len(batch[1]) > 0:
                for result in zip(batch[1], max(self._network(batch[2]), 1)[1]):
                    self._updater.update_image(
                        Image(result[0], Label(result[1].item()))
                    )
        self._drain.drain_duration(Duration.total, start)


class Reader:
    __slots__ = ['_engine', '_drain']

    def __init__(self, engine: Engine, drain: Drain):
        self._engine = engine
        self._drain = drain

    def read_urls(self) -> List[str]:
        start = time()
        with self._engine.connect() as connection:
            urls = [
                u[0] for u in connection.execute(
                    '''
                    select url
                    from images
                    where kind = 'photo'
                      and label = 'unknown'
                    '''
                )
            ]
        self._drain.drain_duration(Duration.reading, start)
        return urls


# TODO: shorten training code in notebook and use Network, Gallery in jupyter.
class Network(Module):
    __slots__ = ['_sequential']

    def __init__(self):
        super().__init__()
        self._sequential = Sequential(
            Conv2d(3, 3, 5, padding=2),
            ReLU(),
            MaxPool2d(2),
            Dropout(0.3),
            View(-1, 213900),
            Linear(213900, 5)
        )

    def forward(self, x: Tensor) -> Tensor:
        return self._sequential(x)


class View(Module):
    __slots__ = ['_shape']

    def __init__(self, *shape: int):
        super().__init__()
        self._shape = shape

    def forward(self, x: Tensor) -> Tensor:
        return x.view(*self._shape)


class Updater:
    __slots__ = ['_engine', '_drain']

    def __init__(self, engine: Engine, drain: Drain):
        self._engine = engine
        self._drain = drain

    def update_image(self, image: Image):
        start = time()
        with self._engine.connect() as connection:
            connection.execute(
                'update images set label = %s where url = %s',
                image.label.name,
                image.url
            )
        self._drain.drain_duration(Duration.update, start)
        self._drain.drain_number(Number(image.label.value))


class Gallery(Dataset):
    __slots__ = ['_urls', '_session', '_timeout', '_transforms', '_drain']

    def __init__(
        self,
        urls: List[str],
        session: Session,
        timeout: float,
        drain: Drain
    ):
        self._urls = urls
        self._session = session
        self._timeout = timeout
        self._transforms = Compose(
            [
                ToTensor(),
                Resize((460, 620)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )
        self._drain = drain

    def __getitem__(self, index: int) -> Tuple[str, Tensor]:
        start = time()
        try:
            response = self._session.get(
                self._urls[index],
                timeout=self._timeout,
                headers={'User-Agent': 'RampartBot/0.0.1'}
            )
        except RequestException:
            _logger.exception(
                'Gallery failed to read the image',
                extra={'url': self._urls[index]}
            )
            return self._urls[index], empty(0)
        finally:
            self._drain.drain_duration(Duration.loading, start)
        if response.status_code != codes.ok:
            _logger.error(
                'Gallery got non-ok status',
                extra={'url': self._urls[index], 'code': response.status_code}
            )
            return self._urls[index], empty(0)
        image = open(BytesIO(response.content))
        if image.mode == 'RGBA':
            canvas = new('RGBA', image.size, 'white')
            canvas.paste(image, (0, 0), image)
            image = canvas.convert('RGB')
        return self._urls[index], self._transforms(image)

    def __len__(self) -> int:
        return len(self._urls)


def _collate(batch: List[Tuple[str, Tensor]]) -> Tuple[List[str], List[str], Tensor]:
    urls, pairs = [], []
    for pair in batch:
        if pair[1].size()[0] == 0:
            urls.append(pair[0])
        else:
            pairs.append(pair)
    if len(pairs) == 0:
        return urls, [], empty(0)
    bundle = default_collate(pairs)
    return urls, bundle[0], bundle[1]
