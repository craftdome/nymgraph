package client

import (
	"database/sql"
)

type Client interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
	Ping() error
	Close() error
}
