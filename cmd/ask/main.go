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
	question := flag.String("q", "", "user question to preprocess and embed")
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

	result, err := engine.ProcessQuery(ctx, *question)
	if err != nil {
		log.Fatal("query processing failed", zap.Error(err))
	}

	fmt.Printf("Question:          %s\n", result.Question)
	fmt.Printf("Intent:            %s\n", result.Intent)
	fmt.Printf("Policy filter:     %s\n", result.PolicyTypeFilter)
	fmt.Printf("Embed model:       %s\n", result.EmbedModel)
	fmt.Printf("Vector dimensions: %d\n", result.VectorDimensions)
}
