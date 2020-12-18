from io import BytesIO
from typing import List, Tuple
from PIL.Image import open
from requests import Session, codes
from requests.exceptions import RequestException
from sqlalchemy.engine.base import Engine
from torch.nn import (
    Module, Sequential, ReLU, Conv2d, MaxPool2d, Dropout, Linear
)
from torch import Tensor, empty
from torch.utils.data.dataloader import default_collate
from torch.utils.data.dataset import Dataset
from torchvision.transforms import Compose, ToTensor, Resize, Normalize
from rampart.config import GalleryConfig
from rampart.logging import get_logger
from rampart.models import Image

_logger = get_logger(__name__)


class Recognizer(Module):
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


class Gallery(Dataset):
    __slots__ = ['_timeout', '_session', '_urls', '_transforms']

    def __init__(
        self,
        config: GalleryConfig,
        session: Session,
        engine: Engine
    ):
        self._timeout = config.timeout
        self._session = session
        with engine.connect() as connection:
            proxy = connection.execute(
                '''
                select url
                from images
                where kind = 'photo'
                  and label = 'unknown'
                '''
            )
            self._urls: List[str] = [u[0] for u in proxy]
        self._transforms = Compose(
            [
                ToTensor(),
                Resize((460, 620)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )

    def __getitem__(self, index: int) -> Tuple[str, Tensor]:
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
        if response.status_code != codes.ok:
            _logger.error(
                'Gallery got non-ok status',
                extra={'url': self._urls[index], 'code': response.status_code}
            )
            return self._urls[index], empty(0)
        return (
            self._urls[index],
            self._transforms(open(BytesIO(response.content)))
        )

    def __len__(self) -> int:
        return len(self._urls)


def collate(batch: List[Tuple[str, Tensor]]):
    bundle = [p for p in batch if p[1].size()[0] != 0]
    return bundle if len(bundle) == 0 else default_collate(bundle)


# TODO: change parsing time into recognition_time. Consider simultaneous flat/
#  image insert, so we need to memorise just classification.
class Storer:
    __slots__ = ['_engine']

    def __init__(self, engine: Engine):
        self._engine = engine

    def store_image(self, image: Image):
        with self._engine.connect() as connection:
            connection.execute(
                'update images set label = %s where url = %s',
                image.label.name,
                image.url
            )
