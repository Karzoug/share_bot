package db

const createQuery string = `CREATE TABLE IF NOT EXISTS "expenses" (
	"id"	INTEGER,
	"borrower_id"	INTEGER NOT NULL,
	"lender_id"	INTEGER NOT NULL,
	"sum"	INTEGER DEFAULT 0,
	"request_id"	INTEGER NOT NULL,
	"returned"	INTEGER DEFAULT 0,
	"approved"	INTEGER DEFAULT 0,
	FOREIGN KEY("borrower_id") REFERENCES "users"("id"),
	FOREIGN KEY("request_id") REFERENCES "requests"("id"),
	FOREIGN KEY("lender_id") REFERENCES "users"("id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "requests" (
	"id"	INTEGER,
	"comment"	TEXT,
	"date"	INTEGER NOT NULL,
	"chat_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER,
	"username"	TEXT NOT NULL UNIQUE,
	"chat_id"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE INDEX IF NOT EXISTS expenses_not_returned_index ON expenses (borrower_id) WHERE returned = 0;
CREATE INDEX IF NOT EXISTS users_username_index ON users (username);`
