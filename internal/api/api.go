package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"

	"github.com/Karzoug/share_bot/internal/api/model"
)

const telegramAPIUrlFmt = "https://api.telegram.org/bot%s/%s"

type API struct {
	cfg         Config
	BotName     string
	retryClient *retryablehttp.Client
	logger      *zap.Logger
}

func New(ctx context.Context, cfg Config, logger *zap.Logger) (API, error) {
	const op = "api: new"

	a := API{
		cfg:         cfg,
		logger:      logger,
		retryClient: retryablehttp.NewClient(),
	}
	a.retryClient.RetryMax = 5
	a.retryClient.Logger = nil // disable (warn: sensitive data could be logged)

	_, botName, err := a.GetMe(ctx)
	if err != nil {
		return API{}, fmt.Errorf("%s: %w", op, err)
	}
	a.BotName = botName
	return a, nil
}

func (a API) GetMe(ctx context.Context) (uint64, string, error) {
	const op = "api: get me"

	req, err := retryablehttp.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf(telegramAPIUrlFmt, a.cfg.Token, "getMe"),
		http.NoBody)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}
	resp, err := a.retryClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	var p fastjson.Parser
	v, err := p.ParseBytes(body)
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}
	return v.GetUint64("id"), string(v.GetStringBytes("username")), nil
}

func (a API) DeleteMessage(ctx context.Context, chatID, messageID int64) (bool, error) {
	const op = "api: delete message"

	req := model.DeleteMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	httpReq, err := retryablehttp.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf(telegramAPIUrlFmt, a.cfg.Token, "deleteMessage"),
		reqBytes)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.retryClient.Do(httpReq)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	return true, nil
}

func (a API) SendMessage(ctx context.Context, msg string, chatID int64, rm model.ReplyMarkup) (uint64, error) {
	const op = "api: send message"

	req := model.SendMessageRequest{
		ChatID:      chatID,
		Text:        msg,
		ParseMode:   "HTML",
		ReplyMarkup: rm,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	httpReq, err := retryablehttp.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf(telegramAPIUrlFmt, a.cfg.Token, "sendMessage"),
		reqBytes)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.retryClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("%s: unexpected status code %d", op, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	var p fastjson.Parser
	v, err := p.ParseBytes(respBody)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return v.GetUint64("message_id"), nil
}
