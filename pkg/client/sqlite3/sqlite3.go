package sqlite3

import (
	"database/sql"
	"fmt"
	"github.com/craftdome/nymgraph/pkg/client"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

type Config struct {
	DBFileName string
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"file:%s",
		c.DBFileName,
	)
}

func (c Config) String() string {
	return c.DSN()
}

func NewClient(cfg Config) (client.Client, error) {
	db, err := sql.Open("sqlite", cfg.DSN())
	if err != nil {
		return nil, errors.Wrapf(err, "sql.Open(\"sqlite\", %s)", cfg.DSN())
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db.Ping()")
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	return db, nil
}
