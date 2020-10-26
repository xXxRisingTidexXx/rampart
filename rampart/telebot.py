from logging import INFO, basicConfig
from telegram.ext import Updater, CommandHandler, Filters, MessageHandler
from rampart.config import get_config
from rampart.handlers import TelebotHandler

# TODO: add JSON logging.
if __name__ == '__main__':
    basicConfig(level=INFO)
    config = get_config().telebot
    handler = TelebotHandler(config.handler)
    updater = Updater(config.token)
    updater.dispatcher.add_handler(CommandHandler('start', handler.get_start))
    updater.dispatcher.add_handler(CommandHandler('help', handler.get_help))
    filters = Filters.regex(config.pattern)
    updater.dispatcher.add_handler(MessageHandler(filters, handler.get_search))
    updater.dispatcher.add_handler(MessageHandler(~filters, handler.get_confusion))
    updater.start_polling()
    updater.idle()
