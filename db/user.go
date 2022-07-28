package db

import (
	"database/sql"
	"fmt"
	"log"
)

// InitUser calls on start bot-user interaction in private chat
func InitUser(user User) {
	mustOpenDb()

	dbUser, ok := getUser(user.Username)
	if !ok {
		addUser(&user)
	} else {
		updateUser(dbUser.Id, &user)
	}
}

// getUser tries to find user in database by username
func getUser(username string) (dbUser User, exist bool) {
	err := db.QueryRow("SELECT id, username, telegram_id FROM users WHERE username = $1", username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.TelegramId)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false
		} else {
			log.Fatal(fmt.Errorf("get User query error: %w", err))
		}
	}

	return dbUser, true
}

// addUser inserts user to database and add id to user
func addUser(user *User) {
	stmt, err := db.Prepare("INSERT INTO users(username, telegram_id) VALUES($1, $2)")
	if err != nil {
		log.Fatal(fmt.Errorf("add User prepare query error: %w", err))
	}
	res, err := stmt.Exec(user.Username, user.TelegramId)
	if err != nil {
		log.Fatal(fmt.Errorf("add User query execution error: %w", err))
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(fmt.Errorf("add User get last id error: %w", err))
	}
	user.Id = int(lastId)
}

// updateUser updates only TelegramId for User with Id = id
func updateUser(id int, new *User) {
	stmt, err := db.Prepare("UPDATE users SET telegram_id = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(fmt.Errorf("update User prepare query error: %w", err))
	}
	_, err = stmt.Exec(new.TelegramId, id)
	if err != nil {
		log.Fatal(fmt.Errorf("update User query execution error: %w", err))
	}

	new.Id = id
}
