from fastapi import FastAPI

app = FastAPI()


@app.get('/')
def _get_root():
    return {'greeting': 'Hello, world!'}
