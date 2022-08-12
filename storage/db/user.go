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

	dbUser, exist := st.GetUserByUsername(user.Username)
	if !exist {
		dbUser.Username = user.Username
		dbUser.ChatId = user.ChatId
		st.addUser(&dbUser)
		return
	}

	dbUser.ChatId = user.ChatId
	st.updateUser(&dbUser)
}

// getUser tries to find user in database by username
func (st Storage) GetUserByUsername(username string) (dbUser storage.User, exist bool) {
	err := st.db.QueryRow("SELECT id, username, chat_id FROM users WHERE username = $1", username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.User{}, false
		} else {
			log.Panic(e.Wrap("can't get user by username", err))
		}
	}

	return dbUser, true
}

func (st Storage) GetUserById(id int) (dbUser storage.User, exist bool) {
	err := st.db.QueryRow("SELECT id, username, chat_id FROM users WHERE id = $1", id).Scan(&dbUser.Id, &dbUser.Username, &dbUser.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.User{}, false
		} else {
			log.Panic(e.Wrap("can't get user by id", err))
		}
	}

	return dbUser, true
}

func (st Storage) IsUserExist(username string) (exist bool) {
	err := st.db.QueryRow("SELECT case when exists (SELECT NULL FROM users WHERE username = $1) then 1 else 0 end", username).Scan(&exist)
	if err != nil {
		exist = false
		if err != sql.ErrNoRows {
			log.Panic(e.Wrap("can't get is user exist", err))
		}
	}

	return
}

func (st Storage) addUser(u *storage.User) {
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
	u.Id = lastId
	return
}

func (st Storage) updateUser(u *storage.User) {
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
