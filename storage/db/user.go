package db

import (
	"database/sql"
	"log"
	"share_bot/lib/e"
	"share_bot/storage"
)

// SaveUser calls on start bot-user interaction in private chat
func (st Storage) SaveUser(user storage.User) {
	st.mustOpenDb()

	dbUser, exist := st.getUser(user.Username)
	if !exist {
		st.addUser(&tableUser{
			Username: user.Username,
			ChatId:   user.ChatId,
		})
		return
	}

	dbUser.ChatId = user.ChatId
	st.updateUser(&dbUser)
}

// getUser tries to find user in database by username
func (st Storage) getUser(username string) (dbUser tableUser, exist bool) {
	err := st.db.QueryRow("SELECT id, username, chat_id FROM users WHERE username = $1", username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return tableUser{Username: username}, false
		} else {
			log.Panic(e.Wrap("can't get user", err))
		}
	}

	return dbUser, true
}

func (st Storage) addUser(u *tableUser) {
	var err error
	defer func() {
		if err != nil {
			log.Panic(e.Wrap("can't add user", err))
		}
	}()

	stmt, err := st.db.Prepare("INSERT INTO users(username, chat_id) VALUES($1, $2)")
	if err != nil {
		return
	}
	res, err := stmt.Exec(u.Username, u.ChatId)
	if err != nil {
		return
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return
	}
	u.Id = int(lastId)
	return
}

func (st Storage) updateUser(u *tableUser) {
	var err error
	defer func() {
		if err != nil {
			log.Panic(e.Wrap("can't update user", err))
		}
	}()

	stmt, err := st.db.Prepare("UPDATE users SET username = $1, chat_id = $2 WHERE id = $3")
	if err != nil {
		return
	}
	_, err = stmt.Exec(u.Username, u.ChatId, u.Id)
	if err != nil {
		return
	}
}
