package api

import (
	"context"

	"hrpolicyassistant/internal/query"
	"hrpolicyassistant/internal/retrieval"
)

// AskResult is the output of the full runtime query pipeline (Steps 8–10).
type AskResult struct {
	Query     *query.Result
	Retrieval *retrieval.Result
}

// Asker runs the full runtime query pipeline.
type Asker interface {
	Ask(ctx context.Context, question string) (*AskResult, error)
}
