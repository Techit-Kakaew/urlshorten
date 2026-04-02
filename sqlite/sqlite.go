package sqlite

import (
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Save(shortenURL, originalURL string) error {
	_, err := s.db.Exec("INSERT INTO urls (key, original_url) VALUES (?, ?)", shortenURL, originalURL)
	return err
}

func (s *Storage) OriginalURL(shortenURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE key= ?", shortenURL).Scan(&originalURL)

	return originalURL, err
}
