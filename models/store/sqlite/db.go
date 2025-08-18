package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const create_users_table = `CREATE TABLE IF NOT EXISTS clusters (
    id TEXT PRIMARY KEY,
	email TEXT NOT NULL,
    name TEXT NOT NULL
);`

const create_user = `INSERT INTO users (id, email, name) VALUES (?, ?, ?);`

const user_exists = `SELECT * from users WHERE id = ?;`

type SqliteModel struct {
	dsn string
}

func getDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSqliteModel(dsn string) (*SqliteModel, error) {
	db, err := getDBConnection(dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(create_users_table)
	if err != nil {
		return nil, err
	}

	return &SqliteModel{dsn: dsn}, nil
}

func (m *SqliteModel) Close() {

}
