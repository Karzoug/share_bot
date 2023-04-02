package bot

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"share_bot/internal/logger"
	"share_bot/internal/storage"

	"github.com/NicoNex/echotron/v3"
	"go.uber.org/zap"
)

const setWebhookUrl string = "https://api.telegram.org/bot%s/setWebhook"

type dispatcher struct {
	token string
	*echotron.Dispatcher
}

func NewDispatcher(token string, storage storage.Storage) *dispatcher {
	if token == "" {
		logger.Logger.Fatal("telegram token does not exist")
	}
	newBotFn := func(chatID int64) echotron.Bot {
		return &bot{
			chatID,
			echotron.NewAPI(token),
			storage,
		}
	}
	return &dispatcher{token: token, Dispatcher: echotron.NewDispatcher(token, newBotFn)}
}

// ListenWebhook sets a webhook and listens for incoming updates.
// The addr should be provided in the following format: '<hostname>:<port>/<path>',
// eg: 'https://example.com:443/bot_token'.
// ListenWebhook will then proceed to communicate the webhook url '<hostname>/<path>' to Telegram
// and run a webserver that listens to ':<port>' and handles the path.
// If certificatePath is not empty that certificate file sends to Telegram
func (d dispatcher) ListenWebhook(addr, certificatePath string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}
	if len(certificatePath) == 0 {
		return d.Dispatcher.ListenWebhook(addr)
	}
	whURL := fmt.Sprintf("%s%s", u.Hostname(), u.EscapedPath())
	d.setWebhook(whURL, certificatePath)

	http.HandleFunc(u.EscapedPath(), d.HandleWebhook)
	return http.ListenAndServe(fmt.Sprintf(":%s", u.Port()), nil)
}

func (d dispatcher) setWebhook(whURL, certificatePath string) {
	mustExistsFile(certificatePath)

	contType, reader, err := createReqBody(whURL, certificatePath)
	if err != nil {
		logger.Logger.Fatal("setwebhook with certificate request send error:", zap.Error(err))
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(setWebhookUrl, d.token), reader)
	req.Header.Add("Content-Type", contType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Fatal("setwebhook with certificate request send error:", zap.Error(err))
	}
	resp.Body.Close()
}

func mustExistsFile(filePath string) {
	if _, err := os.Stat(filePath); err != nil {
		logger.Logger.Fatal("can't open file", zap.Error(err))
	}
}

func createReqBody(wURL, filePath string) (string, io.Reader, error) {
	var err error

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	f, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	p2w, _ := bw.CreateFormField("url")
	p2w.Write([]byte(wURL))

	// file part1
	_, fileName := filepath.Split(filePath)
	fw1, _ := bw.CreateFormFile("certificate", fileName)
	io.Copy(fw1, f)

	bw.Close() //write the tail boundry
	return bw.FormDataContentType(), buf, nil
}
