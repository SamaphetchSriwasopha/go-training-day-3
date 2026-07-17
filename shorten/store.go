package shorten

import (
	"context"
	"database/sql"
	"fmt"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) Save(ctx context.Context, shortenURL, rawURL string) error {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO links (code, url) VALUES (?, ?)`, shortenURL, rawURL); err != nil {
		return fmt.Errorf("save shorten: %w", err)
	}
	return nil
}
