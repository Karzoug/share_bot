package user

import (
	"context"
	"fmt"
	"time"

	"github.com/vorlif/spreak"
	"go.uber.org/zap"

	"github.com/Karzoug/share_bot/internal/api"
	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	"github.com/Karzoug/share_bot/internal/model"
	"github.com/Karzoug/share_bot/internal/usecase"
	"github.com/Karzoug/share_bot/internal/usecase/user/repo"
)

var defaultBackgroundTaskTimeout = 3 * time.Second

type Service struct {
	api    api.API
	repo   repo.Repo
	t      *spreak.Localizer
	logger *zap.Logger
}

func New(api api.API,
	userRepo repo.Repo,
	localizer *spreak.Localizer,
	logger *zap.Logger) Service {
	return Service{
		api:    api,
		repo:   userRepo,
		t:      localizer,
		logger: logger,
	}
}

func (s Service) Start(ctx context.Context, user model.User) (string, error) {
	const op = "service: start"

	const helloMsg string = `Привет, %s!
Я бот, который поможет тебе и твоим знакомым не забыть об общих тратах друг друга.
	 
• Чтобы добавить долг, отправь сообщение следующего вида:
<blockquote>/add @username sum comment</blockquote>
В случае, если указанную сумму нужно разделить поровну на нескольких:
<blockquote>/add @username1 @username2 @username3 sum comment</blockquote>
В случае, если долг имеет один комментарий, но разные суммы:
<blockquote>/add @username1 sum1 @usehname2 sum2 comment</blockquote>

Примеры:
<blockquote>/add @anna 800 вино</blockquote>
<blockquote>/add @anna 500 @james 300 вино</blockquote>
<blockquote>/add @viktor @vasya @petya 1700 шашлык</blockquote>
	
• Чтобы быстро узнать, кто и сколько должен, воспользуйся кнопками снизу.

• Чтобы сообщить о возврате, используй появившиеся кнопки в списке твоих долгов.

Начнем?`

	if err := s.repo.Save(ctx, user); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var msg string
	if user.FirstName != "" {
		msg = s.t.Getf(helloMsg, user.FirstName)
	} else {
		msg = s.t.Getf(helloMsg, user.Username)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		_, err := s.api.SendMessage(ctx, msg, user.ID, &apiModel.ReplyKeyboardMarkup{
			Keyboard: [][]apiModel.KeyboardButton{{
				{Text: s.t.Get(usecase.GetUserDebtsText)},
				{Text: s.t.Get(usecase.GetDebtsOwedToUserText)},
			}},
			IsPersistent:   true,
			ResizeKeyboard: true,
		})
		if err != nil {
			s.logger.Error("error", zap.String("error message", err.Error()))
		}
	}()

	return "", nil
}
