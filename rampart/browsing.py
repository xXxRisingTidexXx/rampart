from flask import Flask

if __name__ == '__main__':
    app = Flask('rampart.browsing')

    app.run('0.0.0.0', 9211, load_dotenv=False, use_reloader=False)
