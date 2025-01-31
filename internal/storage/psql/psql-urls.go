package psql

import (
	"bybit-parser/internal/lig/logger/sl"
	"bybit-parser/internal/storage"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) SaveURL(url string, alias string) (int64, error) {
	const op = "storage.psql.SaveURL"
	lastInsertId := int64(0)
	err := s.db.QueryRow(`INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id`, url, alias).Scan(&lastInsertId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Printf("%s\t%s\t%s\n", lastInsertId, url, alias)

	return lastInsertId, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.psql.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, sl.Err(err))
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	res, err := s.db.Exec("DELETE FROM url WHERE alias=$1", alias)

	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 0 {
				return fmt.Errorf("Ничего не удалено")
			}
			/* check count and return true/false */
		}
	}

	return nil
}
