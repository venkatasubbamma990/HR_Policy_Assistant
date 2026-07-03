package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"hrpolicyassistant/internal/query"
)

// Handler serves runtime query endpoints.
type Handler struct {
	pipeline *query.Pipeline
	log      *zap.Logger
}

// NewHandler creates an HTTP handler for the query pipeline.
func NewHandler(pipeline *query.Pipeline, log *zap.Logger) *Handler {
	return &Handler{
		pipeline: pipeline,
		log:      log.Named("api"),
	}
}

type queryRequest struct {
	Question string `json:"question"`
}

type queryResponse struct {
	Question         string         `json:"question"`
	Normalized       string         `json:"normalized"`
	Intent           string         `json:"intent"`
	PolicyTypeFilter string         `json:"policy_type_filter"`
	MetadataFilter   map[string]any `json:"metadata_filter,omitempty"`
	EmbedModel       string         `json:"embed_model"`
	VectorDimensions int            `json:"vector_dimensions"`
}

// HandleQuery processes Step 8 and Step 9 for an incoming user question.
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

	result, err := h.pipeline.Process(r.Context(), req.Question)
	if err != nil {
		h.log.Error("query processing failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := queryResponse{
		Question:         result.Question,
		Normalized:       result.Normalized,
		Intent:           result.Intent,
		PolicyTypeFilter: result.PolicyTypeFilter,
		MetadataFilter:   result.MetadataFilter,
		EmbedModel:       result.EmbedModel,
		VectorDimensions: result.VectorDimensions,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// HandleHealth reports service health.
func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
