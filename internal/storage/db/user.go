package db

import (
	"database/sql"
	"share_bot/internal/storage"
	"share_bot/pkg/e"
)

// SaveUser calls on start bot-user interaction in private chat
func (st Storage) SaveUser(user storage.User) (err error) {
	dbUser, exist, err := st.GetUserByUsername(user.Username)
	if err != nil {
		return
	}
	if !exist {
		dbUser.Username = user.Username
		dbUser.ChatId = user.ChatId
		return st.addUser(&dbUser)
	}
	dbUser.ChatId = user.ChatId
	return st.updateUser(&dbUser)
}

// getUser tries to find user in database by username
func (st Storage) GetUserByUsername(username string) (dbUser storage.User, exist bool, err error) {
	st.mustOpenDb()
	err = st.db.QueryRow("SELECT id, username, chat_id FROM users WHERE username = $1", username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.ChatId)
	if err != nil {
		if err != sql.ErrNoRows {
			return storage.User{}, false, e.Wrap("can't get user by username", err)
		}
		return storage.User{}, false, nil
	}
	return dbUser, true, nil
}

func (st Storage) GetUserById(id int) (dbUser storage.User, exist bool, err error) {
	st.mustOpenDb()
	err = st.db.QueryRow("SELECT id, username, chat_id FROM users WHERE id = $1", id).Scan(&dbUser.Id, &dbUser.Username, &dbUser.ChatId)
	if err != nil {
		if err != sql.ErrNoRows {
			return storage.User{}, false, e.Wrap("can't get user by id", err)
		}
		return storage.User{}, false, nil
	}
	return dbUser, true, nil
}

func (st Storage) IsUserExist(username string) (exist bool, err error) {
	err = st.db.QueryRow("SELECT case when exists (SELECT NULL FROM users WHERE username = $1) then 1 else 0 end", username).Scan(&exist)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, e.Wrap("can't get is user exist", err)
		}
		return false, nil
	}
	return
}

func (st Storage) addUser(u *storage.User) (err error) {
	st.mustOpenDb()
	defer func() {
		if err != nil {
			err = e.Wrap("can't add user", err)
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

func (st Storage) updateUser(u *storage.User) (err error) {
	st.mustOpenDb()
	defer func() {
		if err != nil {
			err = e.Wrap("can't update user", err)
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
	return
}
