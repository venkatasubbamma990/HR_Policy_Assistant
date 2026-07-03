package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/vectorstores/pgvector"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/api"
	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/index"
	"hrpolicyassistant/internal/query"
)

// Engine orchestrates the RAG query pipeline and HTTP API.
type Engine struct {
	cfg              *config.Config
	store            pgvector.Store
	queryPipeline    *query.Pipeline
	vectorStoreReady bool
	log              *zap.Logger
}

// NewEngine creates a RAG engine connected to pgvector with query preprocessing.
func NewEngine(ctx context.Context, cfg *config.Config, log *zap.Logger) (*Engine, error) {
	engine := &Engine{
		cfg: cfg,
		log: log.Named("rag"),
	}

	if !cfg.IndexingEnabled() {
		return engine, nil
	}

	store, embedder, err := index.OpenStore(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open vector store: %w", err)
	}

	engine.store = store
	engine.queryPipeline = query.NewPipeline(cfg, embedder, log)
	engine.vectorStoreReady = true
	return engine, nil
}

// ProcessQuery runs Step 8 (preprocess) and Step 9 (embed question).
func (e *Engine) ProcessQuery(ctx context.Context, question string) (*query.Result, error) {
	if e.queryPipeline == nil {
		return nil, fmt.Errorf("query pipeline not configured; set DATABASE_URL and OPENAI_API_KEY")
	}
	return e.queryPipeline.Process(ctx, question)
}

// Run starts the HTTP query API until shutdown.
func (e *Engine) Run(ctx context.Context) error {
	e.log.Info("RAG engine started",
		zap.String("embed_model", e.cfg.EmbedModel),
		zap.String("chat_model", e.cfg.ChatModel),
		zap.Bool("openai_configured", e.cfg.OpenAIAPIKey != ""),
		zap.Bool("vector_store_ready", e.vectorStoreReady),
		zap.String("http_addr", e.cfg.HTTPAddr()),
	)

	if !e.vectorStoreReady {
		e.log.Warn("query API disabled; configure DATABASE_URL and OPENAI_API_KEY")
		<-ctx.Done()
		return ctx.Err()
	}

	handler := api.NewHandler(e.queryPipeline, e.log)
	server := api.NewServer(e.cfg.HTTPAddr(), handler, e.log)

	err := server.Start(ctx)

	if e.vectorStoreReady {
		_ = e.store.Close()
	}

	if err != nil && err != context.Canceled {
		return err
	}

	e.log.Info("RAG engine shutting down")
	return err
}

// Close releases vector store resources.
func (e *Engine) Close() error {
	if e.vectorStoreReady {
		return e.store.Close()
	}
	return nil
}

// Store returns the pgvector store for retrieval (Step 10).
func (e *Engine) Store() pgvector.Store {
	return e.store
}
