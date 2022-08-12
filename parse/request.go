package parse

import (
	"errors"
	"regexp"
	"share_bot/storage"
	"strconv"
	"strings"
)

// AddMessage parses add expense message from user
func AddMessage(message string) (exps []storage.Expense, comment string, e error) {
	message = strings.TrimSpace(message)
	re := regexp.MustCompile(`(@[0-9a-z_]+) (\d+) ([^@]*)`)
	str := re.FindAllStringSubmatch(message, -1)

	sum := 0
	var err error
	exps = make([]storage.Expense, 0, len(str))
	for i := 0; i < len(str); i++ {
		sum, err = strconv.Atoi(str[i][2])
		if err != nil {
			return nil, "", errors.New("parse sum error")
		}
		exps = append(exps, storage.Expense{
			Borrower: strings.TrimPrefix(str[i][1], "@"),
			Sum:      sum,
		})
	}

	if len(exps) > 0 {
		comment = strings.TrimSpace(str[len(str)-1][3])
	}

	return
}
