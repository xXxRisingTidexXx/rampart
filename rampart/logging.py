from logging import Logger, getLogger, StreamHandler, Handler
from pythonjsonlogger.jsonlogger import JsonFormatter


# TODO: add permanent label "app".
def get_logger(name: str) -> Logger:
    logger = getLogger(name)
    logger.addHandler(get_handler())
    return logger


def get_handler() -> Handler:
    handler = StreamHandler()
    handler.setFormatter(JsonFormatter())
    return handler
