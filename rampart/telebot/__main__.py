from logging import INFO, basicConfig
from telegram.ext import Updater, CommandHandler, Filters, MessageHandler
from rampart.config import get_config
from rampart.telebot.serve import Server

# TODO: add JSON logging.
if __name__ == '__main__':
    basicConfig(level=INFO)
    config = get_config().telebot
    server = Server(config.server)
    updater = Updater(config.token)
    updater.dispatcher.add_handler(
        CommandHandler(config.start_command, server.get_start)
    )
    updater.dispatcher.add_handler(
        CommandHandler(config.help_command, server.get_help)
    )
    filters = Filters.regex(config.search_pattern)
    updater.dispatcher.add_handler(MessageHandler(filters, server.get_search))
    updater.dispatcher.add_handler(MessageHandler(~filters, server.get_confusion))
    updater.start_polling()
    updater.idle()
