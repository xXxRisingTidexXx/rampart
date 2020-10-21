from logging import basicConfig, INFO, getLogger
from os import getenv
from telegram import Update
from telegram.ext import Updater, CommandHandler, CallbackContext
from rampart.search import search_flats, Query

_logger = getLogger('rampart.telebot')
_FLOORS = {'-': 0, 'низько': 1, 'високо': 2}
_ROOM_NUMBERS = {'-': 0, 'одна': 1, 'дві': 2, 'три': 3, 'багато': 4}


# TODO: leverage optuna to set the hyperparameters.
# TODO: add JSON logging.
def _main():
    updater = Updater(getenv('RAMPART_TELEBOT_TOKEN'))
    updater.dispatcher.add_handler(CommandHandler('start', _get_start))
    updater.dispatcher.add_handler(CommandHandler('search', _get_search))
    updater.dispatcher.add_handler(CommandHandler('help', _get_help))
    updater.start_polling()
    updater.idle()


def _get_start(update: Update, _: CallbackContext):
    update.message.reply_markdown_v2(
        'Привіт\\! Дозволь ввести тебе в хід справ: я \\- щось по типу гугла, але для ву'
        'зької категорії речей :\\)\\. Та результати знаходжу так само \\- за текстовими'
        ' запитами; щоправда, в них доволі специфічний формат, менше схожий на природню '
        'мову\\. Кожен запит має вигляд:```\n\n/search <місто> <ціна> <поверх> <кімнати>'
        '```\n\nІєрогліфи в трикутних дужках \\- деякі шаблонні значення, котрі тобі вар'
        'то замінити на власні\\. Якщо не знаєш чи не хочеш вказувати якийсь із пунктів,'
        ' то постав `\\-` \\. Наприклад:```\n\n/search Київ 75000 високо \\-```\n\n Дост'
        'упні значення:\n\n\\- місто \\- актуальна назва довільного міста України\n\\- ц'
        'іна \\- доступна для тебе сума в USD\n\\- поверх \\- одне з двох значень: `висо'
        'ко` чи `низько`\n\\- кімнати \\- одне зі значень: `одна` , `дві` , `три` чи `ба'
        'гато`\n\nP\\.S\\. Я ігнорую всі некомандні повідомлення, так усім жити легше\\.'
    )


def _get_search(update: Update, context: CallbackContext):
    if len(context.args) != 4:
        update.message.reply_text(
            f'Невірна кількість параметрів - ти надав {len(context.args)}, а треба 4.'
        )
        return
    city = context.args[0]
    if city == '-':
        city = 'Київ'
    price = _float(context.args[1])
    if price < 0:
        update.message.reply_text('Перевір, будь ласка, ціну - вона некоректна.')
        return
    floor = _FLOORS.get(context.args[2], -1)
    if floor < 0:
        update.message.reply_text(
            'Хм, я не зрозумів отриманий ідентифікатор поверху (підглянь у /start).'
        )
        return
    room_number = _ROOM_NUMBERS.get(context.args[3], -1)
    if room_number < 0:
        update.message.reply_text(
            'Теекс, це якесь дивне число кімнат. Перевір введення (звірся зі /start).'
        )
        return
    flats = search_flats(Query(city, price, floor, room_number))
    if len(flats) == 0:
        update.message.reply_text('На жаль, мені нічого не вдалося знайти.')


def _float(value: str) -> float:
    if value == '-':
        return 0
    try:
        return float(value)
    except ValueError:
        return -1


def _get_help(update: Update, _: CallbackContext):
    update.message.reply_text(
        'До твоїх послуг доступні такі команди:\n\n/start - довідка щодо формату пошуков'
        'ого запиту\n/search - пошук житла\n/help - це повідомлення\n'
    )


if __name__ == '__main__':
    basicConfig(level=INFO)
    _main()
