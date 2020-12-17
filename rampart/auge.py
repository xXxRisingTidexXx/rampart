from requests.adapters import HTTPAdapter
from sqlalchemy import create_engine
from torch import load
from rampart.config import get_config
from rampart.recognition import Recognizer, Gallery
from requests import Session


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
    gallery = Gallery(create_engine(config.auge.dsn), session)

    # recognizer = Recognizer()
    # recognizer.load_state_dict(load(config.auge.model_path))
    session.close()


if __name__ == '__main__':
    _main()
