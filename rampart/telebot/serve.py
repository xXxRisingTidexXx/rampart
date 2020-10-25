from os import getenv
from telegram import Update
from telegram.ext import (
    Updater, CommandHandler, CallbackContext, MessageHandler, Filters
)
from rampart.telebot.search import Query, Searcher, Housing


def serve():
    server = Server()
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'))
    updater.dispatcher.add_handler(CommandHandler('start', server.get_start))
    updater.dispatcher.add_handler(CommandHandler('help', server.get_help))
    filters = Filters.regex(
        r'^\s*(\S+) +([1-9]\d*|0|-) +(-|низько|високо) +(-|одна|дві|три|багато)\s*$'
    )
    updater.dispatcher.add_handler(MessageHandler(filters, server.get_flats))
    updater.dispatcher.add_handler(MessageHandler(~filters, server.get_confusion))
    updater.start_polling()
    updater.idle()


class Server:
    __slots__ = [
        '_searcher',
        '_start_template',
        '_help_template',
        '_nothing_template',
        '_flat_template',
        '_confusion_template'
    ]
    _any = '-'
    _cities = {_any: 'Київ'}
    _floors = {_any: 0, 'низько': 1, 'високо': 2}
    _room_numbers = {_any: 0, 'одна': 1, 'дві': 2, 'три': 3, 'багато': 4}
    _housings = {Housing.primary: 'первинка', Housing.secondary: 'вторинка'}

    def __init__(self):
        self._searcher = Searcher()
        with open('templates/start.html') as stream:
            self._start_template = stream.read()
        with open('templates/help.html') as stream:
            self._help_template = stream.read()
        with open('templates/nothing.html') as stream:
            self._nothing_template = stream.read()
        with open('templates/flat.html') as stream:
            self._flat_template = stream.read()
        with open('templates/confusion.html') as stream:
            self._confusion_template = stream.read()

    def get_start(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._start_template)

    def get_help(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._help_template)

    def get_flats(self, update: Update, context: CallbackContext):
        groups = context.match.groups()
        flats = self._searcher.search_flats(
            Query(
                self._cities.get(groups[0], groups[0]),
                0 if groups[1] == self._any else float(groups[1]),
                self._floors[groups[2]],
                self._room_numbers[groups[3]]
            )
        )
        if len(flats) == 0:
            update.message.reply_html(self._nothing_template)
        for flat in flats:
            update.message.reply_html(
                self._flat_template.format(
                    flat.url,
                    flat.address,
                    self._housings[flat.housing],
                    flat.price,
                    flat.room_number,
                    flat.total_area,
                    flat.total_floor,
                    flat.floor
                )
            )

    def get_confusion(self, update: Update, _: CallbackContext):
        if update.message:
            update.message.reply_html(self._confusion_template)
