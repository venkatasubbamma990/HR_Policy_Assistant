package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/index"
	"hrpolicyassistant/internal/ingest"
	"hrpolicyassistant/internal/logger"
	"hrpolicyassistant/internal/rag"
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

	log = log.Named("hrpolicy")
	log.Info("starting HR Policy Assistant",
		zap.String("documents_dir", cfg.DocumentsDir),
		zap.String("log_level", cfg.LogLevel),
		zap.String("log_format", cfg.LogFormat),
		zap.Int("chunk_min_tokens", cfg.ChunkMinTokens),
		zap.Int("chunk_max_tokens", cfg.ChunkMaxTokens),
		zap.Int("chunk_overlap_tokens", cfg.ChunkOverlapTokens),
		zap.Bool("indexing_enabled", cfg.IndexingEnabled()),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pipeline := ingest.NewPipeline(cfg, log)
	result, err := pipeline.Run()
	if err != nil {
		log.Fatal("document ingestion failed", zap.Error(err))
	}

	if cfg.IndexingEnabled() {
		indexer := index.NewIndexer(cfg, log)
		indexResult, err := indexer.Run(ctx, result.Documents, result.Chunks)
		if err != nil {
			log.Fatal("vector indexing failed", zap.Error(err))
		}
		if indexResult.Skipped {
			log.Info("using existing vector index")
		}
	} else {
		log.Warn("vector indexing disabled; set DATABASE_URL and OPENAI_API_KEY to enable pgvector")
	}

	engine, err := rag.NewEngine(ctx, cfg, log)
	if err != nil {
		log.Fatal("rag engine init failed", zap.Error(err))
	}

	if err := engine.Run(ctx); err != nil && err != context.Canceled {
		log.Fatal("rag engine stopped with error", zap.Error(err))
	}

	log.Info("HR Policy Assistant stopped",
		zap.Int("documents", len(result.Documents)),
		zap.Int("chunks", len(result.Chunks)),
	)
}
