package main

import (
	"context"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/index"
	"hrpolicyassistant/internal/ingest"
	"hrpolicyassistant/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("config: " + err.Error())
	}

	log, err := logger.New(logger.Options{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	if err != nil {
		panic("logger: " + err.Error())
	}
	defer func() { _ = log.Sync() }()

	log = log.Named("index")
	if !cfg.IndexingEnabled() {
		log.Fatal("vector indexing requires DATABASE_URL and OPENAI_API_KEY")
	}

	log.Info("starting policy indexing pipeline",
		zap.String("documents_dir", cfg.DocumentsDir),
		zap.String("database_url", cfg.DatabaseURL),
		zap.String("embed_model", cfg.EmbedModel),
		zap.Int("embed_dimensions", cfg.EmbedDimensions),
		zap.String("collection", cfg.PGVectorCollection),
	)

	ctx := context.Background()

	pipeline := ingest.NewPipeline(cfg, log)
	result, err := pipeline.Run()
	if err != nil {
		log.Fatal("document ingestion failed", zap.Error(err))
	}

	indexer := index.NewIndexer(cfg, log)
	indexResult, err := indexer.Run(ctx, result.Documents, result.Chunks)
	if err != nil {
		log.Fatal("vector indexing failed", zap.Error(err))
	}

	if indexResult.Skipped {
		log.Info("indexing skipped; policies unchanged",
			zap.Int("chunk_count", indexResult.ChunkCount),
		)
		return
	}

	log.Info("indexing finished",
		zap.Int("document_count", len(result.Documents)),
		zap.Int("chunk_count", indexResult.ChunkCount),
		zap.String("content_hash", indexResult.ContentHash),
	)
}
