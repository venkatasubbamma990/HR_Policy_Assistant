package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/index"
)

// Engine orchestrates the RAG pipeline for HR policy Q&A.
type Engine struct {
	cfg              *config.Config
	store            pgvector.Store
	vectorStoreReady bool
	log              *zap.Logger
}

// NewEngine creates a RAG engine connected to the pgvector store when configured.
func NewEngine(ctx context.Context, cfg *config.Config, log *zap.Logger) (*Engine, error) {
	engine := &Engine{
		cfg: cfg,
		log: log.Named("rag"),
	}

	if !cfg.IndexingEnabled() {
		return engine, nil
	}

	store, _, err := index.OpenStore(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open vector store: %w", err)
	}

	engine.store = store
	engine.vectorStoreReady = true
	return engine, nil
}

// Run starts the assistant loop.
// Query retrieval and LLM generation will be added in subsequent steps.
func (e *Engine) Run(ctx context.Context) error {
	e.log.Info("RAG engine started",
		zap.String("embed_model", e.cfg.EmbedModel),
		zap.String("chat_model", e.cfg.ChatModel),
		zap.Bool("openai_configured", e.cfg.OpenAIAPIKey != ""),
		zap.Bool("vector_store_ready", e.vectorStoreReady),
	)
	e.log.Info("query handling not yet implemented; waiting for shutdown signal")

	<-ctx.Done()

	if e.vectorStoreReady {
		_ = e.store.Close()
	}

	e.log.Info("RAG engine shutting down", zap.Error(ctx.Err()))
	return ctx.Err()
}

// Store returns the pgvector store for retrieval (query-time use).
func (e *Engine) Store() pgvector.Store {
	return e.store
}
