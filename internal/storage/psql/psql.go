package psql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.psql.New"

	db, err := sql.Open(
		"postgres", dbPath,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//тут должна быть миграция, но пока так
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
		id SERIAL NOT NULL,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB Ping failed: %w", err)
	}

	stid, err := db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if _, err := stid.Exec(); err != nil {
		fmt.Println("idx")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, err
}
