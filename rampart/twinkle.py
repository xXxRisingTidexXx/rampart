from os import getenv
from sqlalchemy import create_engine
from rampart.logging import get_logger
from rampart.classification import Classifier, Reader, Query, RoomNumber, Floor

_logger = get_logger('rampart.twinkle')


def _main():
    engine = create_engine(getenv('RAMPART_DSN'))
    reader = Reader(engine)
    classifier = Classifier()
    flats = reader.read_flats(Query('Київ', 0, RoomNumber.one, Floor.any))
    if len(flats) > 0:
        flats['is_relevant'] = classifier.classify_flats(flats.drop(columns=['id']))
    else:
        flats['is_relevant'] = []
    print(len(flats[flats['is_relevant']]))
    engine.dispose()


if __name__ == '__main__':
    _main()
