from sys import exit
from uuid import uuid4
from pandas import read_csv, DataFrame
from rampart.logging import get_logger

_logger = get_logger('rampart.coquus')


# TODO: setup basic logger with lineno & timestamp.
def _main():
    flats = read_csv('scientific/twinkle.csv').drop(columns=['url'])
    group_number = flats['group'].nunique()
    groups = ['training', 'validation', 'testing']
    if group_number != len(groups):
        _logger.critical(
            'Coquus got invalid group number',
            extra={'actual': group_number, 'expected': len(groups)}
        )
        exit(1)
    tag = uuid4().hex
    for i, group in enumerate(groups):
        _augment(flats[flats['group'] == i].drop(columns=['group']), tag, group)


def _augment(flats: DataFrame, tag: str, group: str):
    print(flats['query'].unique())


if __name__ == '__main__':
    _main()
