from enum import Enum, unique
from io import BytesIO
from typing import List
from PIL.Image import open, new, Image as Picture
from requests import Session, codes
from sqlalchemy.engine.base import Engine
from torch import load, no_grad
from torchvision.transforms import Compose, ToTensor, Resize, Normalize
from torchvision.models import Inception3
from rampart.config import LoaderConfig
from rampart.logging import get_logger

_logger = get_logger('rampart.recognition')


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
    __slots__ = ['_session', '_timeout', '_user_agent']

    def __init__(self, config: LoaderConfig, session: Session):
        self._session = session
        self._timeout = config.timeout
        self._user_agent = config.user_agent

    def load_image(self, url: str) -> 'Image':
        response = self._session.get(
            url,
            timeout=self._timeout,
            headers={'User-Agent': self._user_agent}
        )
        if response.status_code != codes.ok:
            raise RuntimeError(f'Loader got non-ok status {response.status_code}')
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
    __slots__ = ['_transforms', '_network']

    def __init__(self, path: str):
        self._transforms = Compose(
            [
                ToTensor(),
                Resize((299, 299)),
                Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
            ]
        )
        self._network = Inception3(5, init_weights=False)
        self._network.load_state_dict(load(path))
        self._network.eval()

    @no_grad()
    def recognize_image(self, image: Image) -> Image:
        pass


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
