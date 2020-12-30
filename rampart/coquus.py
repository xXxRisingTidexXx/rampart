from sys import exit
from uuid import uuid4
from pandas import read_csv, DataFrame, concat
from numpy import ndarray
from numpy.random import rand
from lightgbm import Dataset
from rampart.config import get_config
from rampart.logging import get_logger

_logger = get_logger('rampart.coquus')


# TODO: setup basic logger with lineno & timestamp.
def _main():
    config = get_config()
    flats = read_csv(config.coquus.input_path).drop(columns=['url'])
    group_number = flats['group'].nunique()
    groups = ['training', 'validation', 'testing']
    if group_number != len(groups):
        _logger.critical(
            'Coquus got invalid group number',
            extra={'actual': group_number, 'expected': len(groups)}
        )
        exit(1)
    datasets = {group: _serialize(flats, i) for i, group in enumerate(groups)}
    datasets['validation'].set_reference(datasets['training'])
    tag = uuid4().hex
    for group, dataset in datasets.items():
        dataset.save_binary(config.coquus.output_format.format(tag, group))
        dataset.save_binary(
            config.coquus.output_format.format('latest', group)
        )


def _serialize(flats: DataFrame, i: int) -> Dataset:
    origin = flats[flats['group'] == i].drop(columns=['group'])
    merge = concat(
        [origin] + [_augment(origin, i) for i in range(1, 32)],
        ignore_index=True
    )
    return Dataset(
        merge.drop(columns=['relevance', 'query']),
        merge['relevance'],
        group=merge.groupby(['query']).size(),
        categorical_feature=[
            'desired_room_number',
            'desired_floor',
            'housing'
        ],
        silent=True
    )


def _augment(origin: DataFrame, i: int):
    copy = origin.copy()
    copy['actual_price'] *= _shake(0.001, len(origin))
    copy['total_area'] *= _shake(0.005, len(origin))
    copy['living_area'] *= _shake(0.01, len(origin))
    copy['kitchen_area'] *= _shake(0.01, len(origin))
    copy['ssf'] *= _shake(0.007, len(origin))
    copy['izf'] *= _shake(0.007, len(origin))
    copy['gzf'] *= _shake(0.007, len(origin))
    copy['query'] += i * origin['query'].nunique()
    return copy


def _shake(deviation: float, size: int) -> ndarray:
    return 1 - deviation + 2 * deviation * rand(size)


if __name__ == '__main__':
    _main()
