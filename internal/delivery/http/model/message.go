package model

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Karzoug/share_bot/internal/model"
)

var (
	reMsg                   = regexp.MustCompile(`((?:@[0-9a-zA-Z_]+ )+)(\d+) ([^@]*)`)
	ErrInvalidMessageFormat = errors.New("invalid message format")
)

type Message struct {
	Text string
	ID   int64
	Date int64
}

func ParseDebts(msg Message) ([]model.Debt, error) {
	msg.Text = strings.TrimSpace(msg.Text)
	str := reMsg.FindAllStringSubmatch(msg.Text, -1)

	comment := strings.TrimSpace(str[len(str)-1][3])

	var (
		sum int64
		err error
	)
	debts := make([]model.Debt, len(str))
	for i := 0; i < len(str); i++ {
		sum, err = strconv.ParseInt(str[i][2], 10, 64)
		if err != nil {
			return nil, ErrInvalidMessageFormat
		}
		names := strings.Fields(str[i][1])
		if len(names) > 1 {
			sum /= int64(len(names))
		}
		for _, v := range names {
			debts[i] = model.Debt{
				DebtorUsername: strings.TrimPrefix(v, "@"),
				Sum:            sum,
				Comment:        comment,
				Date:           time.Unix(msg.Date, 0),
			}
		}
	}
	return debts, nil
}
