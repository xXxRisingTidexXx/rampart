from flask import Flask

app = Flask('rampart.browsing')


@app.route('/')
def index():
    return 'Hello, world!'
