from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from torch import load, no_grad, max
from torch.utils.data.dataloader import DataLoader
from rampart.config import get_config
from rampart.models import Image
from rampart.recognition import Recognizer, Gallery, collate, Storer
from requests import Session


# TODO: read more about docker --ipc flag used to satisfy multiprocessing.
# TODO: shorten training code in notebook and use Recognizer, Gallery in jupyter.
# TODO: add label for abandoned images (404/unavailable).
@no_grad()
def _main():
    config = get_config()
    session = Session()
    session.mount(
        'https://',
        HTTPAdapter(
            pool_maxsize=config.auge.thread_number,
            max_retries=config.auge.retry_limit
        )
    )
    engine = create_engine(config.auge.dsn)
    loader = DataLoader(
        Gallery(config.auge.gallery, session, engine),
        config.auge.thread_number,
        num_workers=config.auge.thread_number,
        collate_fn=collate
    )
    storer = Storer(engine)
    recognizer = Recognizer()
    recognizer.load_state_dict(load(config.auge.model_path))
    recognizer.eval()
    for batch in loader:
        if len(batch) == 2:
            for result in zip(batch[0], max(recognizer(batch[1]), 1)[1]):
                storer.store_image(Image(result[0], result[1].item()))
    session.close()
    engine.dispose()


if __name__ == '__main__':
    _main()
