from logging import basicConfig, INFO
from os import getenv
from telegram import Update
from telegram.ext import (
    Updater, CommandHandler, CallbackContext, MessageHandler, Filters
)


# TODO: leverage optuna to set the hyperparameters.
# TODO: add JSON logging.
def _main():
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'))
    updater.dispatcher.add_handler(CommandHandler('start', _get_start))
    updater.dispatcher.add_handler(CommandHandler('r', _get_r))
    updater.dispatcher.add_handler(CommandHandler('help', _get_help))
    updater.dispatcher.add_handler(MessageHandler(~Filters.command, _get_default))
    updater.start_polling()
    updater.idle()


def _get_start(update: Update, _: CallbackContext):
    update.message.reply_markdown_v2(
        'Привіт\\! Дозволь ввести тебе в хід справ: я \\- щось по типу гугла, але для ву'
        'зької категорії речей :\\)\\. Та результати знаходжу так само \\- за текстовими'
        ' запитами; щоправда, в них доволі специфічний формат, менше схожий на природню '
        'мову\\. Кожен запит має вигляд:```\n\n/r <місто> <ціна> <поверх> <кімнати>```\n'
        '\nІєрогліфи в трикутних дужках \\- деякі шаблонні значення, котрі тобі варто за'
        'мінити на власні\\. Якщо не знаєш чи не хочеш вказувати якийсь із пунктів, то п'
        'остав `\\-` \\. Наприклад:```\n\n/r Київ 75000 високо \\-```\n\n Доступні значе'
        'ння:\n\n\\- місто \\- актуальна назва довільного міста України\n\\- ціна \\- до'
        'ступна для тебе сума в USD\n\\- поверх \\- одне з двох значень: `високо` чи `ни'
        'зько`\n\\- кімнати \\- одне зі значень: `одна` , `дві` , `три` чи `багато`\n'
    )


def _get_r(update: Update, context: CallbackContext):
    update.message.reply_text('Пошук')


def _get_help(update: Update, _: CallbackContext):
    update.message.reply_text(
        'До твоїх послуг доступні такі команди:\n\n/start - довідка щодо формату пошуков'
        'ого запиту\n/r - пошук житла\n/help - це повідомлення\n'
    )


def _get_default(update: Update, _: CallbackContext):
    update.message.reply_text('Вибач, звісно, але я не зрозумів, що ти маєш на увазі\n')


if __name__ == '__main__':
    basicConfig(level=INFO)
    _main()
