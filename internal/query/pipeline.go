package query

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/llm"
)

// Result is a preprocessed and embedded user question ready for retrieval.
type Result struct {
	Question         string
	Normalized       string
	SearchText       string
	Intent           string
	PolicyTypeFilter string
	MetadataFilter   map[string]any
	Vector           []float32
	EmbedModel       string
	VectorDimensions int
}

// Pipeline handles runtime question preprocessing and embedding.
type Pipeline struct {
	cfg      *config.Config
	embedder embeddings.Embedder
	log      *zap.Logger
}

// NewPipeline creates a query pipeline using the same embedder model as indexing.
func NewPipeline(cfg *config.Config, embedder embeddings.Embedder, log *zap.Logger) *Pipeline {
	if log == nil {
		log = zap.NewNop()
	}
	return &Pipeline{
		cfg:      cfg,
		embedder: embedder,
		log:      log.Named("query"),
	}
}

// Process executes Step 8 (preprocess) and Step 9 (embed question).
func (p *Pipeline) Process(ctx context.Context, question string) (*Result, error) {
	normalized := Normalize(question)
	if normalized == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	searchText := NormalizedForSearch(normalized)
	intent := DetectIntent(searchText)
	filter := MetadataFilter(intent)

	p.log.Info("received user question",
		zap.String("question", normalized),
		zap.String("intent", intent),
		zap.Any("metadata_filter", filter),
	)

	vector, err := p.embedder.EmbedQuery(ctx, normalized)
	if err != nil {
		err = llm.WrapAPIReachabilityError(p.cfg, err)
		p.log.Error("failed to embed question", zap.Error(err))
		return nil, fmt.Errorf("embed question: %w", err)
	}

	result := &Result{
		Question:         normalized,
		Normalized:       normalized,
		SearchText:       searchText,
		Intent:           intent,
		PolicyTypeFilter: intent,
		MetadataFilter:   filter,
		Vector:           vector,
		EmbedModel:       p.cfg.EmbedModel,
		VectorDimensions: len(vector),
	}

	p.log.Info("question embedded",
		zap.String("embed_model", result.EmbedModel),
		zap.Int("vector_dimensions", result.VectorDimensions),
		zap.String("intent", result.Intent),
	)

	return result, nil
}
