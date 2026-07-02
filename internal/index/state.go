package index

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const indexStateTable = "hr_policy_index_state"

// IndexState tracks the last successful vector index run.
type IndexState struct {
	ContentHash string
	EmbedModel  string
	ChunkCount  int
	IndexedAt   time.Time
}

type stateStore struct {
	conn *pgx.Conn
}

func newStateStore(ctx context.Context, connURL string) (*stateStore, error) {
	conn, err := pgx.Connect(ctx, connURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	store := &stateStore{conn: conn}
	if err := store.ensureTable(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return store, nil
}

func (s *stateStore) Close(ctx context.Context) {
	if s.conn != nil {
		_ = s.conn.Close(ctx)
	}
}

func (s *stateStore) ensureTable(ctx context.Context) error {
	_, err := s.conn.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT PRIMARY KEY,
			content_hash TEXT NOT NULL,
			embed_model TEXT NOT NULL,
			chunk_count INT NOT NULL,
			indexed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`, indexStateTable))
	return err
}

func (s *stateStore) Get(ctx context.Context) (*IndexState, error) {
	row := s.conn.QueryRow(ctx, fmt.Sprintf(`
		SELECT content_hash, embed_model, chunk_count, indexed_at
		FROM %s
		WHERE id = 1`, indexStateTable))

	var state IndexState
	err := row.Scan(&state.ContentHash, &state.EmbedModel, &state.ChunkCount, &state.IndexedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &state, nil
}

func (s *stateStore) Save(ctx context.Context, state IndexState) error {
	_, err := s.conn.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (id, content_hash, embed_model, chunk_count, indexed_at)
		VALUES (1, $1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			content_hash = EXCLUDED.content_hash,
			embed_model = EXCLUDED.embed_model,
			chunk_count = EXCLUDED.chunk_count,
			indexed_at = EXCLUDED.indexed_at`,
		indexStateTable,
	), state.ContentHash, state.EmbedModel, state.ChunkCount, state.IndexedAt)

	return err
}
