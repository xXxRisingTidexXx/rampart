from os import getenv
from telegram import Update
from telegram.ext import Updater, CommandHandler, CallbackContext
from rampart.telebot.search import Query, Searcher


def serve():
    server = Server()
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'))
    updater.dispatcher.add_handler(CommandHandler('start', server.get_start))
    updater.dispatcher.add_handler(CommandHandler('search', server.get_search))
    updater.dispatcher.add_handler(CommandHandler('help', server.get_help))
    updater.start_polling()
    updater.idle()


class Server:
    __slots__ = [
        '_searcher',
        '_start_template',
        '_args_template',
        '_price_template',
        '_floor_template',
        '_room_number_template',
        '_nothing_template',
        '_flat_template',
        '_help_template'
    ]
    _any = '-'
    _floors = {_any: 0, 'низько': 1, 'високо': 2}
    _room_numbers = {_any: 0, 'одна': 1, 'дві': 2, 'три': 3, 'багато': 4}

    def __init__(self):
        self._searcher = Searcher()
        with open('templates/start.html') as stream:
            self._start_template = stream.read()
        with open('templates/args.html') as stream:
            self._args_template = stream.read()
        with open('templates/price.html') as stream:
            self._price_template = stream.read()
        with open('templates/floor.html') as stream:
            self._floor_template = stream.read()
        with open('templates/room_number.html') as stream:
            self._room_number_template = stream.read()
        with open('templates/nothing.html') as stream:
            self._nothing_template = stream.read()
        with open('templates/flat.html') as stream:
            self._flat_template = stream.read()
        with open('templates/help.html') as stream:
            self._help_template = stream.read()

    def get_start(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._start_template)

    def get_search(self, update: Update, context: CallbackContext):
        if len(context.args) != 4:
            update.message.reply_html(self._args_template.format(len(context.args)))
            return
        city = context.args[0]
        if city == self._any:
            city = 'Київ'
        price = self._float(context.args[1])
        if price < 0:
            update.message.reply_html(self._price_template.format(context.args[1]))
            return
        floor = self._floors.get(context.args[2], -1)
        if floor < 0:
            update.message.reply_html(self._floor_template)
            return
        room_number = self._room_numbers.get(context.args[3], -1)
        if room_number < 0:
            update.message.reply_html(self._room_number_template)
            return
        flats = self._searcher.search_flats(Query(city, price, floor, room_number))
        if len(flats) == 0:
            update.message.reply_html(self._nothing_template)
            return
        for flat in flats:
            update.message.reply_html(
                self._flat_template.format(flat.url, flat.address)
            )

    def _float(self, value: str) -> float:
        if value == self._any:
            return 0
        try:
            return float(value)
        except ValueError:
            return -1

    def get_help(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._help_template)
