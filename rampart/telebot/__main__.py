# TODO: leverage optuna to set the hyperparameters.
# TODO: add JSON logging.
from logging import INFO, basicConfig
from rampart.telebot.serve import serve

if __name__ == '__main__':
    basicConfig(level=INFO)
    serve()
