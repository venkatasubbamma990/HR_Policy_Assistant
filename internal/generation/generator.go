package generation

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"go.uber.org/zap"

	"hrpolicyassistant/internal/config"
	"hrpolicyassistant/internal/llm"
	"hrpolicyassistant/internal/retrieval"
)

// Generator produces cited answers from retrieved policy chunks.
type Generator struct {
	model llms.Model
	cfg   *config.Config
	log   *zap.Logger
}

// NewGenerator creates an answer generator using the configured chat model.
func NewGenerator(cfg *config.Config, log *zap.Logger) (*Generator, error) {
	model, err := llm.NewChatModel(cfg)
	if err != nil {
		return nil, err
	}

	if log == nil {
		log = zap.NewNop()
	}

	return &Generator{
		model: model,
		cfg:   cfg,
		log:   log.Named("generation"),
	}, nil
}

// Generate creates a cited answer from retrieved chunks (Step 11).
func (g *Generator) Generate(ctx context.Context, question string, retrievalResult *retrieval.Result) (*Result, error) {
	if retrievalResult == nil || len(retrievalResult.Chunks) == 0 {
		g.log.Warn("no retrieved chunks available for answer generation",
			zap.String("question", question),
		)
		return NoContextResult(g.cfg.ChatModel), nil
	}

	userPrompt := BuildUserPrompt(question, retrievalResult.Chunks)
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	g.log.Info("generating cited answer",
		zap.String("question", question),
		zap.String("chat_model", g.cfg.ChatModel),
		zap.Int("context_chunks", len(retrievalResult.Chunks)),
	)

	response, err := g.model.GenerateContent(ctx, messages)
	if err != nil {
		g.log.Error("answer generation failed", zap.Error(err))
		return nil, fmt.Errorf("generate answer: %w", err)
	}
	if response == nil || len(response.Choices) == 0 || strings.TrimSpace(response.Choices[0].Content) == "" {
		return nil, fmt.Errorf("empty response from chat model")
	}

	answerText := strings.TrimSpace(response.Choices[0].Content)
	citations := ExtractCitations(answerText, retrievalResult.Chunks)

	g.log.Info("answer generated",
		zap.Int("citation_count", len(citations)),
		zap.Int("answer_length", len(answerText)),
	)

	return &Result{
		Text:      answerText,
		ChatModel: g.cfg.ChatModel,
		Citations: citations,
	}, nil
}
