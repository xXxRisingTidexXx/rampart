from os import getenv
from fastapi import FastAPI
from sqlalchemy import create_engine

app = FastAPI()
_engine = create_engine(getenv('RAMPART_DSN'))


@app.get('/')
def _read_root():
    return {'greeting': 'Great!'}
