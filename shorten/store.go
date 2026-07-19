package shorten

import (
	"context"
	"database/sql"
	"fmt"
)

type Dstore struct {
	fdDb *sql.DB
}

func NewStore(db *sql.DB) *Dstore {
	return &Dstore{fdDb: db}
}

func (s *Dstore) Save(ctx context.Context, shortenURL, rawURL string) error {
	if _, err := s.fdDb.ExecContext(ctx, `INSERT INTO links (code, url) VALUES (?, ?)`, shortenURL, rawURL); err != nil {
		return fmt.Errorf("save shorten: %w", err)
	}
	return nil
}
