from os import getenv
from telegram import Update
from telegram.ext import Updater, CommandHandler, CallbackContext


def serve():
    server = Server()
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'))
    updater.dispatcher.add_handler(CommandHandler('start', server.get_start))
    # updater.dispatcher.add_handler(CommandHandler('search', _get_search))
    updater.dispatcher.add_handler(CommandHandler('help', server.get_help))
    updater.start_polling()
    updater.idle()


class Server:
    def __init__(self):
        with open('templates/start.html') as stream:
            self._start_template = stream.read()
        with open('templates/help.html') as stream:
            self._help_template = stream.read()

    def get_start(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._start_template)

    def get_help(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._help_template)
