from io import BytesIO
from logging import getLogger
from typing import List, Tuple
from PIL.Image import open
from requests import Session
from sqlalchemy.engine.base import Engine
from torch.nn import (
    Module, Sequential, ReLU, Conv2d, MaxPool2d, Dropout, Linear
)
from torch import Tensor, empty
from torch.utils.data.dataloader import default_collate
from torch.utils.data.dataset import Dataset
from torchvision.transforms import Compose, ToTensor, Resize, Normalize

_logger = getLogger(__name__)


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
    __slots__ = ['_session', '_urls', '_transforms']

    def __init__(self, session: Session, engine: Engine):
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

    def __getitem__(self, index: int) -> Tuple[Tensor, str]:
        response = self._session.get(
            self._urls[index],
            headers={'User-Agent': 'RampartBot/0.0.1'}
        )
        if response.status_code != 200:
            _logger.error('Gallery got non-ok status')
            return empty(0), self._urls[index]
        return (
            self._transforms(open(BytesIO(response.content))),
            self._urls[index]
        )

    def __len__(self) -> int:
        return len(self._urls)


def collate(batch: List[Tuple[Tensor, str]]):
    return default_collate([p for p in batch if p[0].size()[0] != 0])
