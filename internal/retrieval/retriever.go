package retrieval

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/query"
)

// Chunk is a policy chunk returned by similarity search.
type Chunk struct {
	Content     string  `json:"content"`
	Score       float32 `json:"score"`
	Source      string  `json:"source"`
	PolicyType  string  `json:"policy_type"`
	SectionPath string  `json:"section_path"`
	DocumentID  string  `json:"document_id"`
	ChunkID     string  `json:"chunk_id"`
}

// Result holds retrieved chunks for a query.
type Result struct {
	Chunks       []Chunk `json:"chunks"`
	TopK         int     `json:"top_k"`
	FilterUsed   bool    `json:"filter_used"`
	FallbackUsed bool    `json:"fallback_used"`
}

// Retriever performs pgvector similarity search over indexed policy chunks.
type Retriever struct {
	store pgvector.Store
	cfg   *config.Config
	log   *zap.Logger
}

// NewRetriever creates a retriever backed by pgvector.
func NewRetriever(store pgvector.Store, cfg *config.Config, log *zap.Logger) *Retriever {
	if log == nil {
		log = zap.NewNop()
	}
	return &Retriever{
		store: store,
		cfg:   cfg,
		log:   log.Named("retrieval"),
	}
}

// Retrieve finds the top-K most relevant chunks for an embedded query (Step 10).
func (r *Retriever) Retrieve(ctx context.Context, q *query.Result) (*Result, error) {
	if q == nil {
		return nil, fmt.Errorf("query result is required")
	}

	r.log.Info("starting similarity search",
		zap.String("question", q.Question),
		zap.Int("top_k", r.cfg.RetrievalTopK),
		zap.Float32("min_score", r.cfg.RetrievalMinScore),
		zap.Any("metadata_filter", q.MetadataFilter),
	)

	filterUsed := len(q.MetadataFilter) > 0
	docs, err := r.search(ctx, q.Question, q.MetadataFilter)
	if err != nil {
		return nil, err
	}

	fallbackUsed := false
	if len(docs) == 0 && filterUsed {
		r.log.Info("no chunks matched metadata filter; retrying without filter",
			zap.Any("metadata_filter", q.MetadataFilter),
		)
		docs, err = r.search(ctx, q.Question, nil)
		if err != nil {
			return nil, err
		}
		fallbackUsed = true
		filterUsed = false
	}

	chunks := toChunks(docs)
	r.log.Info("similarity search completed",
		zap.Int("result_count", len(chunks)),
		zap.Bool("filter_used", filterUsed),
		zap.Bool("fallback_used", fallbackUsed),
	)

	for _, chunk := range chunks {
		r.log.Debug("retrieved chunk",
			zap.String("chunk_id", chunk.ChunkID),
			zap.String("source", chunk.Source),
			zap.String("section_path", chunk.SectionPath),
			zap.Float32("score", chunk.Score),
		)
	}

	return &Result{
		Chunks:       chunks,
		TopK:         r.cfg.RetrievalTopK,
		FilterUsed:   filterUsed,
		FallbackUsed: fallbackUsed,
	}, nil
}

func (r *Retriever) search(ctx context.Context, question string, filter map[string]any) ([]schema.Document, error) {
	opts := []vectorstores.Option{
		vectorstores.WithScoreThreshold(r.cfg.RetrievalMinScore),
	}
	if len(filter) > 0 {
		opts = append(opts, vectorstores.WithFilters(filter))
	}

	docs, err := r.store.SimilaritySearch(ctx, question, r.cfg.RetrievalTopK, opts...)
	if err != nil {
		return nil, fmt.Errorf("similarity search: %w", err)
	}

	return docs, nil
}

func toChunks(docs []schema.Document) []Chunk {
	chunks := make([]Chunk, 0, len(docs))
	for _, doc := range docs {
		chunks = append(chunks, Chunk{
			Content:     doc.PageContent,
			Score:       doc.Score,
			Source:      metadataString(doc.Metadata, "source"),
			PolicyType:  metadataString(doc.Metadata, "policy_type"),
			SectionPath: metadataString(doc.Metadata, "section_path"),
			DocumentID:  metadataString(doc.Metadata, "document_id"),
			ChunkID:     metadataString(doc.Metadata, "chunk_id"),
		})
	}
	return chunks
}

func metadataString(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	value, ok := metadata[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}
