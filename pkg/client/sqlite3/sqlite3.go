package sqlite3

import (
	"database/sql"
	"fmt"
	"github.com/Tyz3/nymgraph/pkg/client"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
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
	db, err := sql.Open("sqlite3", cfg.DSN())
	if err != nil {
		return nil, errors.Wrapf(err, "sql.Open(\"sqlite3\", %s)", cfg.DSN())
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db.Ping()")
	}

	return db, nil
}
