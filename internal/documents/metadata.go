package documents

import (
	"fmt"
	"regexp"
	"strings"
)

var headerFieldPattern = regexp.MustCompile(`^\*\*(.+?):\*\*\s*(.+)$`)

// Metadata holds document-level fields extracted from policy headers.
type Metadata struct {
	Source        string `json:"source"`
	PolicyType    string `json:"policy_type"`
	DocID         string `json:"doc_id"`
	Version       string `json:"version"`
	EffectiveDate string `json:"effective_date"`
	Company       string `json:"company,omitempty"`
	AppliesTo     string `json:"applies_to,omitempty"`
	LastReviewed  string `json:"last_reviewed,omitempty"`
}

// ParseMetadata extracts header fields from the markdown preamble and filename.
func ParseMetadata(sourceFilename, content string) (Metadata, string, error) {
	title, headerBlock := splitTitleAndHeader(content)

	meta := Metadata{
		Source:     sourceFilename,
		PolicyType: policyTypeFromSource(sourceFilename),
	}

	for _, line := range strings.Split(headerBlock, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line == "---" {
			continue
		}

		matches := headerFieldPattern.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(matches[1]))
		value := strings.TrimSpace(matches[2])

		switch key {
		case "document id":
			meta.DocID = value
		case "version":
			meta.Version = value
		case "effective date":
			meta.EffectiveDate = value
		case "company":
			meta.Company = value
		case "applies to":
			meta.AppliesTo = value
		case "last reviewed":
			meta.LastReviewed = value
		}
	}

	if meta.DocID == "" {
		return Metadata{}, "", fmt.Errorf("document id not found in %q", sourceFilename)
	}

	return meta, title, nil
}

func splitTitleAndHeader(content string) (title, headerBlock string) {
	lines := strings.Split(content, "\n")
	var headerLines []string
	titleFound := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			break
		}

		if !titleFound && strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
			titleFound = true
			continue
		}

		headerLines = append(headerLines, line)
	}

	return title, strings.Join(headerLines, "\n")
}

func policyTypeFromSource(filename string) string {
	base := strings.ToLower(filename)
	base = strings.TrimSuffix(base, ".md")
	base = strings.TrimSuffix(base, "-policy")
	return strings.ReplaceAll(base, "-", "")
}
