from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from torch import load, no_grad, max
from torch.utils.data.dataloader import DataLoader
from rampart.config import get_config
from rampart.models import Label
from rampart.recognition import Recognizer, Gallery
from requests import Session


# TODO: read more about docker --ipc flag used to satisfy multiprocessing.
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
        Gallery(engine, session),
        config.auge.thread_number,
        num_workers=config.auge.thread_number
    )
    recognizer = Recognizer()
    recognizer.load_state_dict(load(config.auge.model_path))
    recognizer.eval()
    for batch in loader:
        for result in zip(max(recognizer(batch[0]), 1)[1], batch[1]):
            print(Label(result[0].item()).name, result[1])
    session.close()
    engine.dispose()


if __name__ == '__main__':
    _main()
