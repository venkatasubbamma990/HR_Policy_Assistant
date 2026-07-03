package query

import "strings"

type intentRule struct {
	policyType string
	keywords   []string
}

// policyIntentRules maps keywords to policy_type metadata used during retrieval.
var policyIntentRules = []intentRule{
	{policyType: "notice", keywords: []string{"notice period", "notice", "resignation notice"}},
	{policyType: "leave", keywords: []string{"privilege leave", "casual leave", "sick leave", "comp off", "take pl", "maternity leave", "paternity leave", "bereavement leave", "marriage leave", "leave", "pl", "cl", "sl"}},
	{policyType: "salary", keywords: []string{"compensation", "reimbursement", "increment", "variable pay", "salary", "ctc", "bonus", "hra", "pay"}},
	{policyType: "wfh", keywords: []string{"work from home", "hybrid work", "remote work", "hybrid", "remote", "wfh"}},
	{policyType: "entry", keywords: []string{"onboarding", "new hire", "joining", "probation", "entry"}},
	{policyType: "exit", keywords: []string{"full and final", "separation", "resignation", "relieving", "exit"}},
}

// DetectIntent infers the most likely policy_type filter from the question.
// Returns an empty string when no confident match is found.
func DetectIntent(searchText string) string {
	searchText = strings.ToLower(strings.TrimSpace(searchText))

	bestType := ""
	bestScore := 0

	for _, rule := range policyIntentRules {
		score := 0
		for _, keyword := range rule.keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword != "" && strings.Contains(searchText, keyword) {
				score += len(keyword)
			}
		}
		if score > bestScore {
			bestScore = score
			bestType = rule.policyType
		}
	}

	return bestType
}

// MetadataFilter builds a pgvector metadata filter from detected intent.
func MetadataFilter(policyType string) map[string]any {
	if policyType == "" {
		return nil
	}
	return map[string]any{
		"policy_type": policyType,
	}
}
