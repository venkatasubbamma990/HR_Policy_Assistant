package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
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
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pipeline := ingest.NewPipeline(cfg, log)
	result, err := pipeline.Run()
	if err != nil {
		log.Fatal("document ingestion failed", zap.Error(err))
	}

	engine := rag.NewEngine(cfg, result.Chunks, log)

	if err := engine.Run(ctx); err != nil && err != context.Canceled {
		log.Fatal("rag engine stopped with error", zap.Error(err))
	}

	log.Info("HR Policy Assistant stopped",
		zap.Int("documents", len(result.Documents)),
		zap.Int("chunks", len(result.Chunks)),
	)
}
