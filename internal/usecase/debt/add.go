package debt

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/xid"
	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	"github.com/Karzoug/share_bot/internal/model"
	"github.com/Karzoug/share_bot/internal/usecase"
)

func (s Service) AddDebts(ctx context.Context,
	debts []model.Debt,
	authorID int64,
	chat model.Chat,
	msgID int64) (string, error) {
	const op = "debt service: add debts"

	const (
		debtAddedMsgFmt               = "@%s сообщил о тратах «%s»\n\n"
		debtAddedSuccessfullyMsg      = "Спасибо! Я запомнила указанные вами траты!"
		needDeletePermissionMsg       = "Бот должен иметь права на удаление сообщений пользователей"
		mentionedUserNotRegisteredMsg = "К сожалению, упомянутый пользователь не зарегистрирован в боте! Может быть порекомендуете меня? 😉"
		needToRegisterMsg             = "Пожалуйста, сначала зарегистрируйтесь в боте 🙏"
		confirmButtonMsg              = "Подтвердить"
	)

	if len(debts) == 0 {
		return "", nil
	}

	// author must be registered
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return "", usecase.NewError(s.t.Get(needToRegisterMsg))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		resDel, err := s.api.DeleteMessage(ctx, chat.ID, msgID)
		if err != nil {
			s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
		}
		if chat.Type != "private" && !resDel {
			_, err := s.api.SendMessage(ctx, s.t.Get(needDeletePermissionMsg), chat.ID, nil)
			if err != nil {
				s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
			}
		}
	}()

	// save debts in repo
	reqID := xid.New().String()
	for i := range debts {
		debts[i].AuthorID = authorID
		debts[i].RequestID = reqID

		if err := s.debtRepo.Save(ctx, debts[i]); err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	// send messages about debts in goroutine
	go func() {
		if len(debts) == 0 {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		if chat.Type == "private" {
			for i := range debts {
				var bld strings.Builder
				if user, err := s.userRepo.GetByUsername(ctx, debts[i].DebtorUsername); err == nil { // if we know the debtor
					bld.WriteString(s.t.Getf(debtAddedMsgFmt, author.Username, debts[i].Comment))
					fmt.Fprintf(&bld, "🧾 @%s: %d %s \n", debts[i].DebtorUsername, debts[i].Sum, s.t.Get("₽"))
					_, err := s.api.SendMessage(ctx,
						bld.String(),
						user.ID,
						apiModel.InlineKeyboardMarkup{
							InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
								Text:         confirmButtonMsg,
								CallbackData: fmt.Sprintf("confirm_debt:%s", reqID),
							}}},
						})
					if err != nil {
						s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
					}
				} else if errors.Is(err, sql.ErrNoRows) || (err == nil && user.ID == 0) { // if we don't know the debtor
					fmt.Fprintf(&bld, "🧾 @%s: %d ₽ «%s»", debts[i].DebtorUsername, debts[i].Sum, debts[i].Comment)
					fmt.Fprint(&bld, "\n\n"+s.t.Get(mentionedUserNotRegisteredMsg))
					_, err := s.api.SendMessage(ctx, bld.String(), chat.ID, nil)
					if err != nil {
						s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
					}
				} else {
					s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
					return
				}
			}
		} else {
			var bld strings.Builder
			bld.WriteString(s.t.Getf(debtAddedMsgFmt, author.Username, debts[0].Comment))
			for i := range debts {
				fmt.Fprintf(&bld, "🧾 @%s: %d %s \n", debts[i].DebtorUsername, debts[i].Sum, s.t.Get("₽"))
			}
			_, err := s.api.SendMessage(ctx,
				bld.String(),
				chat.ID,
				apiModel.InlineKeyboardMarkup{
					InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
						Text:         confirmButtonMsg,
						CallbackData: fmt.Sprintf("confirm_debt:%s", reqID),
					}}},
				})
			if err != nil {
				s.logger.Error("error", zap.String("error message", fmt.Errorf("%s: %w", op, err).Error()))
			}
		}
	}()

	return s.t.Get(debtAddedSuccessfullyMsg), nil
}
