package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/retrieval"
)

// Handler serves runtime query endpoints.
type Handler struct {
	asker Asker
	log   *zap.Logger
}

// NewHandler creates an HTTP handler for the query pipeline.
func NewHandler(asker Asker, log *zap.Logger) *Handler {
	return &Handler{
		asker: asker,
		log:   log.Named("api"),
	}
}

type queryRequest struct {
	Question string `json:"question"`
}

// AskResponse is the API response for a user question.
type AskResponse struct {
	Question         string            `json:"question"`
	Normalized       string            `json:"normalized"`
	Intent           string            `json:"intent"`
	PolicyTypeFilter string            `json:"policy_type_filter"`
	MetadataFilter   map[string]any    `json:"metadata_filter,omitempty"`
	EmbedModel       string            `json:"embed_model"`
	VectorDimensions int               `json:"vector_dimensions"`
	Retrieval        retrievalResponse `json:"retrieval"`
}

type retrievalResponse struct {
	TopK         int               `json:"top_k"`
	FilterUsed   bool              `json:"filter_used"`
	FallbackUsed bool              `json:"fallback_used"`
	Chunks       []retrieval.Chunk `json:"chunks"`
}

// HandleQuery processes Steps 8–10 for an incoming user question.
func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req queryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.asker.Ask(r.Context(), req.Question)
	if err != nil {
		h.log.Error("query processing failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(toAskResponse(result))
}

// HandleHealth reports service health.
func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func toAskResponse(result *AskResult) AskResponse {
	resp := AskResponse{
		Question:         result.Query.Question,
		Normalized:       result.Query.Normalized,
		Intent:           result.Query.Intent,
		PolicyTypeFilter: result.Query.PolicyTypeFilter,
		MetadataFilter:   result.Query.MetadataFilter,
		EmbedModel:       result.Query.EmbedModel,
		VectorDimensions: result.Query.VectorDimensions,
	}

	if result.Retrieval != nil {
		resp.Retrieval = retrievalResponse{
			TopK:         result.Retrieval.TopK,
			FilterUsed:   result.Retrieval.FilterUsed,
			FallbackUsed: result.Retrieval.FallbackUsed,
			Chunks:       result.Retrieval.Chunks,
		}
	}

	return resp
}
