package index

import (
	"context"
	"fmt"
	"time"

	langchainembed "github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/chunk"
	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/documents"
	appembed "hrpolicyassistant/internal/embeddings"
)

// Result describes a vector indexing run.
type Result struct {
	Skipped     bool
	ChunkCount  int
	ContentHash string
	EmbedModel  string
}

// Indexer embeds chunks and stores them in pgvector via langchaingo.
type Indexer struct {
	cfg *config.Config
	log *zap.Logger
}

// NewIndexer creates a vector indexer.
func NewIndexer(cfg *config.Config, log *zap.Logger) *Indexer {
	if log == nil {
		log = zap.NewNop()
	}
	return &Indexer{
		cfg: cfg,
		log: log.Named("index"),
	}
}

// Run embeds and stores chunks when policies changed or force reindex is enabled.
func (idx *Indexer) Run(ctx context.Context, docs []documents.PolicyDocument, chunks []chunk.Chunk) (*Result, error) {
	if !idx.cfg.IndexingEnabled() {
		return nil, fmt.Errorf("indexing requires DATABASE_URL and OPENAI_API_KEY")
	}
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks to index")
	}

	contentHash := Fingerprint(docs)
	stateStore, err := newStateStore(ctx, idx.cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	defer stateStore.Close(ctx)

	current, err := stateStore.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("read index state: %w", err)
	}

	if !idx.cfg.IndexForce && current != nil &&
		current.ContentHash == contentHash &&
		current.EmbedModel == idx.cfg.EmbedModel {
		idx.log.Info("vector index up to date; skipping re-index",
			zap.String("content_hash", contentHash),
			zap.String("embed_model", idx.cfg.EmbedModel),
			zap.Int("chunk_count", current.ChunkCount),
			zap.Time("indexed_at", current.IndexedAt),
		)
		return &Result{
			Skipped:     true,
			ChunkCount:  current.ChunkCount,
			ContentHash: contentHash,
			EmbedModel:  idx.cfg.EmbedModel,
		}, nil
	}

	embedder, err := appembed.NewOpenAIEmbedder(idx.cfg)
	if err != nil {
		return nil, err
	}

	idx.log.Info("starting vector indexing",
		zap.String("embed_model", idx.cfg.EmbedModel),
		zap.Int("embed_dimensions", idx.cfg.EmbedDimensions),
		zap.Int("chunk_count", len(chunks)),
		zap.String("collection", idx.cfg.PGVectorCollection),
		zap.Bool("force", idx.cfg.IndexForce),
	)

	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(idx.cfg.DatabaseURL),
		pgvector.WithEmbedder(embedder),
		pgvector.WithCollectionName(idx.cfg.PGVectorCollection),
		pgvector.WithPreDeleteCollection(true),
		pgvector.WithVectorDimensions(idx.cfg.EmbedDimensions),
		pgvector.WithCollectionMetadata(map[string]any{
			"embed_model":      idx.cfg.EmbedModel,
			"embed_dimensions": idx.cfg.EmbedDimensions,
			"content_hash":     contentHash,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create pgvector store: %w", err)
	}
	defer func() { _ = store.Close() }()

	langDocs := ChunksToDocuments(chunks)
	if err := idx.addDocumentsInBatches(ctx, store, langDocs); err != nil {
		return nil, err
	}

	indexedAt := time.Now().UTC()
	if err := stateStore.Save(ctx, IndexState{
		ContentHash: contentHash,
		EmbedModel:  idx.cfg.EmbedModel,
		ChunkCount:  len(chunks),
		IndexedAt:   indexedAt,
	}); err != nil {
		return nil, fmt.Errorf("save index state: %w", err)
	}

	idx.log.Info("vector indexing completed",
		zap.Int("chunk_count", len(chunks)),
		zap.String("content_hash", contentHash),
		zap.String("embed_model", idx.cfg.EmbedModel),
		zap.Time("indexed_at", indexedAt),
	)

	return &Result{
		Skipped:     false,
		ChunkCount:  len(chunks),
		ContentHash: contentHash,
		EmbedModel:  idx.cfg.EmbedModel,
	}, nil
}

func (idx *Indexer) addDocumentsInBatches(ctx context.Context, store pgvector.Store, docs []schema.Document) error {
	for start := 0; start < len(docs); start += indexBatchSize {
		end := start + indexBatchSize
		if end > len(docs) {
			end = len(docs)
		}

		batch := docs[start:end]
		ids, err := store.AddDocuments(ctx, batch)
		if err != nil {
			return fmt.Errorf("add documents batch %d-%d: %w", start, end, err)
		}

		idx.log.Info("embedded and stored chunk batch",
			zap.Int("batch_start", start),
			zap.Int("batch_end", end),
			zap.Int("batch_size", len(batch)),
			zap.Int("stored_ids", len(ids)),
		)
	}

	return nil
}

// OpenStore creates a pgvector store for query-time retrieval using the same embedder model.
func OpenStore(ctx context.Context, cfg *config.Config) (pgvector.Store, langchainembed.Embedder, error) {
	if !cfg.IndexingEnabled() {
		return pgvector.Store{}, nil, fmt.Errorf("indexing requires DATABASE_URL and OPENAI_API_KEY")
	}

	embedder, err := appembed.NewOpenAIEmbedder(cfg)
	if err != nil {
		return pgvector.Store{}, nil, err
	}

	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(cfg.DatabaseURL),
		pgvector.WithEmbedder(embedder),
		pgvector.WithCollectionName(cfg.PGVectorCollection),
		pgvector.WithVectorDimensions(cfg.EmbedDimensions),
	)
	if err != nil {
		return pgvector.Store{}, nil, fmt.Errorf("open pgvector store: %w", err)
	}

	return store, embedder, nil
}
