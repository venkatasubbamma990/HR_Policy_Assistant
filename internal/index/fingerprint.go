package index

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"hrpolicyassistant/internal/documents"
)

// Fingerprint returns a stable hash of all policy document contents.
func Fingerprint(docs []documents.PolicyDocument) string {
	sorted := append([]documents.PolicyDocument(nil), docs...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Metadata.Source < sorted[j].Metadata.Source
	})

	h := sha256.New()
	for _, doc := range sorted {
		_, _ = io.WriteString(h, doc.Metadata.Source)
		_, _ = io.WriteString(h, "\n")
		_, _ = io.WriteString(h, doc.Raw)
		_, _ = io.WriteString(h, "\n")
	}

	return hex.EncodeToString(h.Sum(nil))
}
