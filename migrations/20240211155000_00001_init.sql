-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "debts" (
	"request_id"         TEXT     NOT NULL,
	"debtor_username"	 TEXT     NOT NULL,
	"author_id"	         INTEGER  NOT NULL,
	"sum"	             INTEGER  NOT NULL,
	"comment"	         TEXT     NOT NULL,
	"date"	             DATETIME NOT NULL,
	"confirmed"	         BOOLEAN  NOT NULL  DEFAULT FALSE,
	PRIMARY KEY("request_id", "debtor_username")
);
CREATE TABLE IF NOT EXISTS "users" (	
	"username"   TEXT    NOT NULL,
    "id"         INTEGER NOT NULL,
	"first_name" TEXT    NOT NULL,
	PRIMARY KEY("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "users_username" ON "users" ("username");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE debts;
DROP TABLE users;
-- +goose StatementEnd
