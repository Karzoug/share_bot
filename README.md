# «SHARE bot»

Простой telegram бот, позволяющий учитывать общие траты в группе знакомых людей и ненавязчиво напоминать о долгах друг другу. 

Проект написан на языке go, в качестве БД используется sqlite. За обертку над Telegram API спасибо проекту <a href="https://github.com/NicoNex/echotron/v3">echotron</a>.

Для взаимодействия с ботом используются команды и кнопки:
* Добавить трату - сообщение-команда в группе с ботом формата:
/add @nickname_друга сумма комментарий_описывающий_трату
* Узнать, кто и сколько должен, - кнопки в приватном общении с ботом
* Вернуть долг - появляющиеся кнопки в списке долгов в приватном общении с ботом.

При присоединении бота к группе:
* при вводе нового долга telegram будет помогать автодополнением имен пользователей;
* работает функционал добавления трат;
* требуются права администратора на удаление сообщений пользователей.

Напоминания о долгах приоритетно будут отправляться ботом должнику в приватном общении, но если на это нет прав, то самому указавшему трату. Это будет случаться только в двух сценариях:
* должник еще не общался с ботом - боты не могут писать пользователям первыми,
* должник не подтвердил этот долг путем нажатия появившейся кнопки.

Настройки, требуемые для запуска бота, задаются в переменных окружения:
* SHARE_BOT_TELEGRAM_TOKEN - секретный токен от BotFather;
* SHARE_BOT_USERNAME - имя пользователя бота.