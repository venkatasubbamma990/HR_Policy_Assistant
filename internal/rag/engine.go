package rag

import (
	"context"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/chunk"
	"hrpolicyassistant/internal/config"
)

// Engine orchestrates the RAG pipeline for HR policy Q&A.
type Engine struct {
	cfg    *config.Config
	chunks []chunk.Chunk
	log    *zap.Logger
}

// NewEngine creates a RAG engine with the given configuration and chunks.
func NewEngine(cfg *config.Config, chunks []chunk.Chunk, log *zap.Logger) *Engine {
	return &Engine{
		cfg:    cfg,
		chunks: chunks,
		log:    log.Named("rag"),
	}
}

// Run starts the interactive assistant loop.
// Embedding, retrieval, and LLM integration will be added in subsequent steps.
func (e *Engine) Run(ctx context.Context) error {
	e.log.Info("RAG engine started",
		zap.Int("indexed_chunks", len(e.chunks)),
		zap.String("embed_model", e.cfg.EmbedModel),
		zap.String("chat_model", e.cfg.ChatModel),
		zap.Bool("openai_configured", e.cfg.OpenAIAPIKey != ""),
	)
	e.log.Info("query handling not yet implemented; waiting for shutdown signal")

	<-ctx.Done()

	e.log.Info("RAG engine shutting down", zap.Error(ctx.Err()))
	return ctx.Err()
}
