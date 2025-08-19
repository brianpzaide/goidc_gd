package sqlite

import (
	"database/sql"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
)

const create_users_table = `CREATE TABLE IF NOT EXISTS clusters (
    id TEXT PRIMARY KEY,
	email TEXT NOT NULL,
    name TEXT NOT NULL
);`

const create_session_index = `CREATE INDEX sessions_expiry_idx ON sessions(expiry);`

const create_user = `INSERT INTO users (id, email, name) VALUES (?, ?, ?);`

const create_sessions_table = `CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);`

const user_exists = `SELECT * from users WHERE id = ?;`

type SqliteModel struct {
	dsn string
}

type SessionManagerApp struct {
	db *sql.DB
	*scs.SessionManager
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

func NewSessionManager(dsn string) (*SessionManagerApp, error) {
	db, err := getDBConnection(dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(create_sessions_table)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(create_session_index)
	if err != nil {
		return nil, err
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	smApp := &SessionManagerApp{
		SessionManager: sessionManager,
	}

	return smApp, nil
}

func (sm *SessionManagerApp) Close() {
	sm.db.Close()
}
