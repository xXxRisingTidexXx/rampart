from os import getenv

from telegram import ReplyKeyboardMarkup, ReplyKeyboardRemove, Update
from telegram.ext import (
    Updater, CommandHandler, MessageHandler, Filters, ConversationHandler,
    CallbackContext
)

GENDER, BIO = range(2)


def start(update: Update, _: CallbackContext) -> int:
    update.message.reply_text(
        'Хай! Я Rampart, маю пару запитань до тебе.\n\nТи хлопець чи дівчина? (якщо волі'
        'єш не відповідати, введи /skip )',
        reply_markup=ReplyKeyboardMarkup([['\U0001F9D4', '\U0001F469']], True, True),
    )
    return GENDER


def set_gender(update: Update, _: CallbackContext) -> int:
    update.message.reply_text(
        'Отримав! А тепер розкажи трішки про себе.',
        reply_markup=ReplyKeyboardRemove(),
    )
    return BIO


def skip_gender(update: Update, _: CallbackContext) -> int:
    update.message.reply_text(
        'Розумію, не питання. Що ж, розкажи мені трішки про себе.',
        reply_markup=ReplyKeyboardRemove(),
    )
    return BIO


def set_bio(update: Update, _: CallbackContext) -> int:
    update.message.reply_text(
        'Дякую за відповідь! Твій голос дуже цінний для нас. Бувай!'
    )
    return ConversationHandler.END


def cancel(update: Update, _: CallbackContext) -> int:
    update.message.reply_text(
        'Прощавай! Ще колись поговоримо :)',
        reply_markup=ReplyKeyboardRemove()
    )
    return ConversationHandler.END


# TODO: leverage optuna to set the hyperparameters.
# TODO: add JSON logging.
def main():
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'), use_context=True)
    updater.dispatcher.add_handler(
        ConversationHandler(
            [CommandHandler('start', start)],
            {
                GENDER: [
                    MessageHandler(
                        Filters.regex('^(\U0001F9D4|\U0001F469)$'),
                        set_gender
                    ),
                    CommandHandler('skip', skip_gender)
                ],
                BIO: [MessageHandler(Filters.text & ~Filters.command, set_bio)]
            },
            [CommandHandler('cancel', cancel)]
        )
    )
    updater.start_polling()
    updater.idle()


if __name__ == '__main__':
    main()
