from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from torch import load, no_grad, max
from torch.utils.data.dataloader import DataLoader
from rampart.config import get_config
from rampart.models import Image, Label
from rampart.recognition import Recognizer, Gallery, collate, Storer
from requests import Session


# TODO: shorten training code in notebook and use Recognizer, Gallery in jupyter.
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
        for url in batch[0]:
            storer.store_image(Image(url, Label.abandoned))
        if len(batch[1]) > 0:
            for result in zip(batch[1], max(recognizer(batch[2]), 1)[1]):
                storer.store_image(Image(result[0], Label(result[1].item())))
    session.close()
    engine.dispose()


if __name__ == '__main__':
    _main()
