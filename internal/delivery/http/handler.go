package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	srvErr "github.com/Karzoug/share_bot/internal/delivery/http/errors"
	httpModel "github.com/Karzoug/share_bot/internal/delivery/http/model"
	"github.com/Karzoug/share_bot/internal/delivery/http/request"
	"github.com/Karzoug/share_bot/internal/delivery/http/response"
	"github.com/Karzoug/share_bot/internal/model"
	"github.com/Karzoug/share_bot/internal/usecase"
	"github.com/Karzoug/share_bot/internal/usecase/debt"
	"github.com/Karzoug/share_bot/internal/usecase/user"
	"github.com/vorlif/spreak"
)

var defaultHandlerTimeout = 3 * time.Second

type handler struct {
	userService user.Service
	debtService debt.Service
	t           *spreak.Localizer
	logger      *zap.Logger
}

func newHandler(userService user.Service,
	debtService debt.Service,
	localizer *spreak.Localizer,
	logger *zap.Logger) handler {
	return handler{
		userService: userService,
		debtService: debtService,
		t:           localizer,
		logger:      logger,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultHandlerTimeout)
	defer cancel()

	const badMessageFormatMsg = "Я не разобралась в описанной вами трате. Проверьте, пожалуйста, синтаксис!"

	var update apiModel.Update
	if err := request.DecodeJSON(w, r, &update); err != nil {
		h.logger.Error("can't decode json", zap.Error(err))
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// route

	if update.Message != nil {
		user := model.User{
			ID:        update.Message.From.ID,
			Username:  update.Message.From.Username,
			FirstName: update.Message.From.FirstName,
		}
		msg := httpModel.Message{
			Text: strings.TrimSpace(update.Message.Text),
			ID:   update.Message.ID,
			Date: update.Message.Date,
		}
		chat := model.Chat{
			Type: update.Message.Chat.Type,
			ID:   update.Message.Chat.ID,
		}

		var (
			method string
			res    string
			err    error
		)
		switch {
		case strings.HasPrefix(update.Message.Text, "/add"):
			method = "bot handler: add debts"
			msg.Text = strings.TrimPrefix(msg.Text, "/add")
			msg.Text = msg.Text[strings.IndexByte(msg.Text, ' ')+1:]

			var debts []model.Debt
			debts, err = httpModel.ParseDebts(msg)
			if err != nil {
				err = usecase.NewError(h.t.Get(badMessageFormatMsg))
			} else {
				res, err = h.debtService.AddDebts(ctx, debts, user.ID, chat, update.Message.ID)
			}
		case strings.HasPrefix(update.Message.Text, h.t.Get(usecase.GetDebtsOwedToUserText)):
			method = "bot handler: get debts owed to user"
			res, err = h.debtService.GetDebtsOwedToUser(ctx, user.ID, chat)
		case strings.HasPrefix(update.Message.Text, h.t.Get(usecase.GetUserDebtsText)):
			method = "bot handler: get user debts"
			res, err = h.debtService.GetUserDebts(ctx, user.Username, chat)
		case update.Message.Text == "/start" && update.Message.Chat.Type == "private":
			method = "bot handler: start"
			res, err = h.userService.Start(ctx, user)
		}
		if err != nil {
			var serviceErr usecase.Error
			if errors.As(err, &serviceErr) {
				res = serviceErr.Error()
			} else {
				srvErr.LogError(method, err, false, h.logger)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if res == "" {
			return
		}
		resp := struct {
			Method string `json:"method"`
			apiModel.SendMessageRequest
		}{
			Method: "sendMessage",
			SendMessageRequest: apiModel.SendMessageRequest{
				ChatID:    update.Message.Chat.ID,
				ParseMode: "HTML",
				Text:      res,
			},
		}
		if err := response.JSON(w, http.StatusOK, resp); err != nil {
			h.logger.Error("can't encode json", zap.Error(err))
		}
		return
	}
	if update.CallbackQuery != nil {
		var (
			method string
			res    string
			err    error
		)
		if strings.HasPrefix(update.CallbackQuery.Data, "confirm_debt:") {
			method = "bot handler: confirm debt"
			splitted := strings.Split(update.CallbackQuery.Data, ":")
			if len(splitted) != 2 {
				h.logger.Error("can't split 'confirm_debt' callback query data",
					zap.String("data", update.CallbackQuery.Data))
				return
			}
			res, err = h.debtService.Confirm(ctx, splitted[1], update.CallbackQuery.From.Username)
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "debt_returned_request:") {
			method = "bot handler: request return debt"
			splitted := strings.Split(update.CallbackQuery.Data, ":")
			if len(splitted) != 3 {
				h.logger.Error("can't split 'debt_returned_request' callback query data",
					zap.String("data", update.CallbackQuery.Data))
				return
			}
			reqID := splitted[1]
			debtorUsername := update.CallbackQuery.From.Username
			res, err = h.debtService.RequestReturn(ctx, reqID, debtorUsername)
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "confirm_return_debt:") {
			method = "bot handler: confirm return debt"
			splitted := strings.Split(update.CallbackQuery.Data, ":")
			if len(splitted) != 3 {
				h.logger.Error("can't split 'confirm_return_debt' callback query data",
					zap.String("data", update.CallbackQuery.Data))
				return
			}
			authorID := update.CallbackQuery.From.ID
			reqID := splitted[1]
			debtorUsername := splitted[2]
			res, err = h.debtService.ConfirmReturn(ctx, authorID, reqID, debtorUsername)
		}

		if err != nil {
			var serviceErr usecase.Error
			if errors.As(err, &serviceErr) {
				res = serviceErr.Error()
			} else {
				srvErr.LogError(method, err, false, h.logger)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if res == "" {
			return
		}
		resp := struct {
			Method string `json:"method"`
			apiModel.AnswerCallbackQuery
		}{
			Method: "answerCallbackQuery",
			AnswerCallbackQuery: apiModel.AnswerCallbackQuery{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            res,
			},
		}
		if err := response.JSON(w, http.StatusOK, resp); err != nil {
			h.logger.Error("can't encode json", zap.Error(err))
		}

		return
	}
}
