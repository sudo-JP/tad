package db

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type TadDatabase struct {
	pool *sql.DB
}

func NewTadDB(dataSourceName string) (*TadDatabase, error) {
	pool, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	pool.SetMaxOpenConns(1)

	if err := pool.Ping(); err != nil {
		return nil, err
	}

	return &TadDatabase{pool: pool}, nil
}

func (db *TadDatabase) Close() error {
	if db.pool != nil {
		return db.pool.Close()
	}
	return nil
}

func translateErr(err error) error {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) &&
		(sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique ||
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey) {
		return ErrAlreadyExists
	}
	return err
}

func requireAffected(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
