package query

import "testing"

func TestNormalize(t *testing.T) {
	got := Normalize("  What is the   notice period?  ")
	want := "What is the notice period?"
	if got != want {
		t.Fatalf("Normalize() = %q, want %q", got, want)
	}
}

func TestDetectIntentNoticePeriod(t *testing.T) {
	intent := DetectIntent(NormalizedForSearch("What is the notice period for L3 engineer?"))
	if intent != "notice" {
		t.Fatalf("intent = %q, want notice", intent)
	}
}

func TestDetectIntentLeave(t *testing.T) {
	intent := DetectIntent(NormalizedForSearch("Can I take PL during notice?"))
	if intent != "leave" {
		t.Fatalf("intent = %q, want leave", intent)
	}
}

func TestDetectIntentSalary(t *testing.T) {
	intent := DetectIntent(NormalizedForSearch("What is the annual CTC for L3?"))
	if intent != "salary" {
		t.Fatalf("intent = %q, want salary", intent)
	}
}

func TestDetectIntentUnknown(t *testing.T) {
	intent := DetectIntent(NormalizedForSearch("Tell me about the office cafeteria"))
	if intent != "" {
		t.Fatalf("intent = %q, want empty", intent)
	}
}

func TestMetadataFilter(t *testing.T) {
	filter := MetadataFilter("notice")
	if filter["policy_type"] != "notice" {
		t.Fatalf("unexpected filter: %#v", filter)
	}
	if MetadataFilter("") != nil {
		t.Fatal("expected nil filter for empty intent")
	}
}
