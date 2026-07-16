package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/logger"
	"hrpolicyassistant/internal/rag"
)

func main() {
	question := flag.String("q", "", "user question to preprocess, embed, and retrieve")
	flag.Parse()

	if *question == "" && flag.NArg() > 0 {
		*question = flag.Arg(0)
	}
	if *question == "" {
		fmt.Fprintln(os.Stderr, "usage: ask -q \"What is the notice period for L3 engineer?\"")
		os.Exit(1)
	}

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

	ctx := context.Background()

	engine, err := rag.NewEngine(ctx, cfg, log)
	if err != nil {
		log.Fatal("rag engine init failed", zap.Error(err))
	}
	defer func() { _ = engine.Close() }()

	result, err := engine.Ask(ctx, *question)
	if err != nil {
		log.Fatal("query failed", zap.Error(err))
	}

	fmt.Printf("Question:      %s\n", result.Query.Question)
	fmt.Printf("Intent:        %s\n", result.Query.Intent)
	fmt.Printf("Embed model:   %s\n", result.Query.EmbedModel)
	fmt.Printf("Chunks found:  %d\n", len(result.Retrieval.Chunks))

	for i, chunk := range result.Retrieval.Chunks {
		fmt.Printf("\n--- Chunk %d (score=%.3f) ---\n", i+1, chunk.Score)
		fmt.Printf("Source:  %s\n", chunk.Source)
		fmt.Printf("Section: %s\n", chunk.SectionPath)
		fmt.Printf("%s\n", truncate(chunk.Content, 300))
	}
}

func truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}
