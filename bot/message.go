package bot

const (
	helloMsg string = `Привет, %v!
Я бот, который поможет тебе и твоим знакомым не забыть об общих тратах друг друга.

* Чтобы добавить трату, отправь сообщение следующего вида:
/add @nickname_друга сумма комментарий_описывающий_трату
* Чтобы быстро узнать, кто и сколько должен, воспользуйся кнопками снизу
* Чтобы вернуть долг, используй появившиеся кнопки в списке твоих долгов

Начнем?`

	toWhomButtonText              string = "Кому должен я?"
	whoMeButtonText               string = "Кто должен мне?"
	needDeletePermissionMsg       string = "Бот должен иметь права на удаление сообщений пользователей"
	somethingWrongTryToStartMsg   string = "Что пошло не так ... Возможно, поможет команда /start"
	toWhomNoExpensesMsg           string = "Вы никому не должны 👍"
	whoMeNoExpensesMsg            string = "Вам никто не должен!"
	returnExpenseButtonMsg        string = "Я вернул этот долг"
	approveExpenseMsg             string = "Вы подтвердили долг. Спасибо!"
	approveReturnExpenseButtonMsg string = "Деньги получил"
	approveReturnExpenseMsg       string = "Вы подтвердили, что получили деньги. Спасибо!"
	sendReturnExpenseMsg          string = "Спасибо! Проверяем ..."
	thanksMsg                     string = "Спасибо!"
	notApprovedMsg                string = "Сумма не была подтверждена должником. Я не напоминаю ему о возврате."
	approveButtonMsg              string = "Подтвердить"
	somethingWrongMsg             string = "Что-то пошло не так ... Уже разбираемся!"
	returnMsg                     string = "@%s сообщил, что отдал вам долг за «%s» в размере %d ₽"
	addMsg                        string = "@%s сообщил о тратах «%s»"
	needToRegisterMsg             string = "%s, сначала зарегистрируйтесь в боте 🙏"
	registerBotButtonMsg          string = "Познакомиться с ботом"
)
