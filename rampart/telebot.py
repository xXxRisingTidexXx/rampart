from logging import INFO, basicConfig
from telegram import Update
from telegram.ext import (
    Updater, CommandHandler, Filters, MessageHandler, CallbackContext
)
from rampart.config import get_config, TelebotHandlerConfig
from rampart.search import Searcher, Query


class Handler:
    __slots__ = [
        '_any',
        '_cities',
        '_floors',
        '_room_numbers',
        '_housings',
        '_start_template',
        '_help_template',
        '_nothing_template',
        '_flat_template',
        '_confusion_template',
        '_searcher'
    ]

    def __init__(self, config: TelebotHandlerConfig):
        self._any = config.any
        self._cities = config.cities
        self._floors = config.floors
        self._room_numbers = config.room_numbers
        self._housings = config.housings
        self._start_template = config.start_template
        self._help_template = config.help_template
        self._nothing_template = config.nothing_template
        self._flat_template = config.flat_template
        self._confusion_template = config.confusion_template
        self._searcher = Searcher(config.searcher)

    def get_start(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._start_template)

    def get_help(self, update: Update, _: CallbackContext):
        update.message.reply_html(self._help_template)

    def get_search(self, update: Update, context: CallbackContext):
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


# TODO: add JSON logging.
if __name__ == '__main__':
    basicConfig(level=INFO)
    setup = get_config().telebot
    updater = Updater(setup.token)
    handler = Handler(setup.handler)
    updater.dispatcher.add_handler(CommandHandler('start', handler.get_start))
    updater.dispatcher.add_handler(CommandHandler('help', handler.get_help))
    filters = Filters.regex(setup.pattern)
    updater.dispatcher.add_handler(MessageHandler(filters, handler.get_search))
    updater.dispatcher.add_handler(MessageHandler(~filters, handler.get_confusion))
    updater.start_polling()
    updater.idle()
