from io import BytesIO
from typing import List, Tuple
from PIL.Image import Image, open
from requests import Session
from sqlalchemy.engine.base import Engine
from torch.nn import (
    Module, Sequential, ReLU, Conv2d, MaxPool2d, Dropout, Linear
)
from torch import Tensor
from torch.utils.data.dataset import Dataset
from torchvision.transforms import Compose, ToTensor, Resize, Normalize


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


class View(Module):
    __slots__ = ['_shape']

    def __init__(self, *shape: int):
        super().__init__()
        self._shape = shape

    def forward(self, x: Tensor):
        return x.view(*self._shape)


class Gallery(Dataset):
    __slots__ = ['_urls', '_transforms']

    def __init__(self, engine: Engine, session: Session):
        with engine.connect() as connection:
            self._urls: List[str] = [
                u[0] for u in connection.execute('select url from images')
            ]
        self._transforms = Compose(
            [
                Download(session),
                ToTensor(),
                Resize((460, 620)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )

    def __len__(self) -> int:
        return len(self._urls)

    def __getitem__(self, index: int) -> Tuple[Tensor, str]:
        return self._transforms(self._urls[index]), self._urls[index]


class Download:
    __slots__ = ['_session']

    def __init__(self, session: Session):
        self._session = session

    def __call__(self, url: str) -> Image:
        return open(
            BytesIO(
                self._session.get(
                    url,
                    headers={'User-Agent': 'RampartBot/0.0.1'}
                ).content
            )
        )
