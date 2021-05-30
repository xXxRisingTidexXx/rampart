from requests import Session
from rampart.mining import Miner
from rampart.geocoding import Geocoder


def _main():
    miner_session = Session()
    miner_session.headers['User-Agent'] = user_agent = 'RampartBot/1.0.0'
    miner = Miner(miner_session)
    geocoder_session = Session()
    geocoder_session.headers['User-Agent'] = user_agent
    geocoder = Geocoder(geocoder_session)
    try:
        flat = miner.mine_flat()
        if flat:
            print(geocoder.geocode_flat(flat))
    finally:
        geocoder_session.close()
        miner_session.close()


if __name__ == '__main__':
    _main()
