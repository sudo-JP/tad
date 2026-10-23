package db

import "database/sql"

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
