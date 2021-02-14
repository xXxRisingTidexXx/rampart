from enum import Enum, unique
from io import BytesIO
from time import time
from typing import List, Tuple
from PIL.Image import open, new
from requests import Session, codes
from requests.exceptions import RequestException
from sqlalchemy.engine.base import Engine
from torch import Tensor, empty, load, no_grad, max
from torch.utils.data.dataloader import default_collate, DataLoader
from torch.utils.data.dataset import Dataset
from torchvision.transforms import Compose, ToTensor, Resize, Normalize
from torchvision.models import AlexNet
from rampart.config import RecognizerConfig
from rampart.logging import get_logger
from rampart.metrics import Drain, Duration, Number

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
        self._network = AlexNet(5)
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
            Gallery(self._reader.read_urls(), self._session, self._timeout),
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
            for span in batch[3]:
                self._drain.drain_duration(Duration.loading, span)
        self._drain.drain_duration(Duration.total, time() - start)
        self._drain.flush()


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
                    select distinct url
                    from images
                    where kind = 'photo'
                      and label = 'unknown'
                    '''
                )
            ]
        self._drain.drain_duration(Duration.reading, time() - start)
        return urls


class Updater:
    __slots__ = ['_engine', '_drain']

    def __init__(self, engine: Engine, drain: Drain):
        self._engine = engine
        self._drain = drain

    def update_image(self, image: 'Image'):
        start = time()
        with self._engine.connect() as connection:
            connection.execute(
                'update images set label = %s where url = %s',
                image.label.name,
                image.url
            )
        self._drain.drain_duration(Duration.update, time() - start)
        self._drain.drain_number(Number(image.label.value))


class Image:
    __slots__ = ['url', 'label']

    def __init__(self, url: str, label: 'Label'):
        self.url = url
        self.label = label


@unique
class Label(Enum):
    unknown = -2
    abandoned = -1
    luxury = 0
    comfort = 1
    junk = 2
    construction = 3
    excess = 4


class Gallery(Dataset):
    __slots__ = ['_urls', '_session', '_timeout', '_transforms']

    def __init__(self, urls: List[str], session: Session, timeout: float):
        self._urls = urls
        self._session = session
        self._timeout = timeout
        self._transforms = Compose(
            [
                ToTensor(),
                Resize((227, 227)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )

    # TODO: use UA from YAML.
    def __getitem__(self, index: int) -> Tuple[str, Tensor, float]:
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
            return self._urls[index], empty(0), time() - start
        span = time() - start
        if response.status_code != codes.ok:
            _logger.error(
                'Gallery got non-ok status',
                extra={'url': self._urls[index], 'code': response.status_code}
            )
            return self._urls[index], empty(0), span
        image = open(BytesIO(response.content))
        if image.mode == 'RGBA':
            canvas = new('RGBA', image.size, 'white')
            canvas.paste(image, (0, 0), image)
            image = canvas.convert('RGB')
        return self._urls[index], self._transforms(image), span

    def __len__(self) -> int:
        return len(self._urls)


def _collate(
    batch: List[Tuple[str, Tensor, float]]
) -> Tuple[List[str], List[str], Tensor, List[float]]:
    bad, good, bundle = [], [], []
    for row in batch:
        if row[1].size()[0] == 0:
            bad.append(row[0])
        else:
            good.append(row[0])
            bundle.append(row[1])
    return (
        bad,
        good,
        empty(0) if len(good) == 0 else default_collate(bundle),
        [r[2] for r in batch]
    )
