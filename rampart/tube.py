from requests import Session
from rampart.gauging import Gauger
from rampart.mining import Miner
from rampart.geocoding import Geocoder


def _main():
    miner_session = Session()
    miner_session.headers['User-Agent'] = user_agent = 'RampartBot/1.0.0'
    miner = Miner(miner_session)
    geocoder_session = Session()
    geocoder_session.headers['User-Agent'] = user_agent
    geocoder = Geocoder(geocoder_session)
    gauger_session = Session()
    gauger_session.headers['User-Agent'] = user_agent
    gauger = Gauger(gauger_session, 'overpass-api.de')
    try:
        flat = miner.mine_flat()
        if flat:
            flat = geocoder.geocode_flat(flat) or flat
            flat = gauger.gauge_flat(flat) or flat
        print(flat)
    finally:
        gauger_session.close()
        geocoder_session.close()
        miner_session.close()


if __name__ == '__main__':
    _main()
