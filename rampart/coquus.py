from uuid import uuid4
from pandas import read_csv, DataFrame, concat
from numpy import ndarray
from numpy.random import rand
from rampart.config import get_config


def _main():
    config = get_config()
    flats = read_csv(config.coquus.input_path).drop(columns=['url'])
    merge = concat(
        [flats] + [_augment(flats) for _ in range(1, 32)],
        ignore_index=True
    )
    merge.to_csv(config.coquus.output_format.format(uuid4().hex), index=False)
    merge.to_csv(config.coquus.output_format.format('latest'), index=False)


def _augment(flats: DataFrame):
    copy = flats.copy()
    copy['actual_price'] *= _noise(0.001, len(flats))
    copy['total_area'] *= _noise(0.005, len(flats))
    copy['living_area'] *= _noise(0.01, len(flats))
    copy['kitchen_area'] *= _noise(0.01, len(flats))
    copy['ssf'] *= _noise(0.007, len(flats))
    copy['izf'] *= _noise(0.007, len(flats))
    copy['gzf'] *= _noise(0.007, len(flats))
    return copy


def _noise(deviation: float, size: int) -> ndarray:
    return 1 - deviation + 2 * deviation * rand(size)


if __name__ == '__main__':
    _main()
