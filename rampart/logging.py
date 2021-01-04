from logging import Logger, getLogger, StreamHandler
from pythonjsonlogger.jsonlogger import JsonFormatter


# TODO: add permanent label "app".
# TODO: setup basic logger with lineno & timestamp.
def get_logger(name: str) -> Logger:
    logger = getLogger(name)
    handler = StreamHandler()
    handler.setFormatter(JsonFormatter())
    logger.addHandler(handler)
    return logger
