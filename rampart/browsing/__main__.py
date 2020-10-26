from rampart.browsing.web import app

if __name__ == '__main__':
    app.run('0.0.0.0', 9211, use_evalex=False, load_dotenv=False)
