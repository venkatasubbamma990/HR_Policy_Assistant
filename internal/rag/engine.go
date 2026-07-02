package rag

import (
	"context"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/documents"
)

// Engine orchestrates the RAG pipeline for HR policy Q&A.
type Engine struct {
	cfg  *config.Config
	docs []documents.PolicyDocument
	log  *zap.Logger
}

// NewEngine creates a RAG engine with the given configuration and documents.
func NewEngine(cfg *config.Config, docs []documents.PolicyDocument, log *zap.Logger) *Engine {
	return &Engine{
		cfg:  cfg,
		docs: docs,
		log:  log.Named("rag"),
	}
}

// Run starts the interactive assistant loop.
// Embedding, retrieval, and LLM integration will be added in subsequent steps.
func (e *Engine) Run(ctx context.Context) error {
	e.log.Info("RAG engine started",
		zap.Int("indexed_documents", len(e.docs)),
		zap.String("embed_model", e.cfg.EmbedModel),
		zap.String("chat_model", e.cfg.ChatModel),
		zap.Bool("openai_configured", e.cfg.OpenAIAPIKey != ""),
	)
	e.log.Info("query handling not yet implemented; waiting for shutdown signal")

	<-ctx.Done()

	e.log.Info("RAG engine shutting down", zap.Error(ctx.Err()))
	return ctx.Err()
}
