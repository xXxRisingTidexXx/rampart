from flask import Flask

if __name__ == '__main__':
    app = Flask('rampart.browsing')

    app.run('0.0.0.0', 9211, True, load_dotenv=False)
