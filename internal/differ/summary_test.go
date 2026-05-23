package differ

import (
	"testing"
)

func makeDriftResult(missing bool, diffs ...FieldDiff) DriftResult {
	return DriftResult{
		Missing: missing,
		Diffs:   diffs,
	}
}

func TestSummarize_AllInSync(t *testing.T) {
	results := []DriftResult{
		makeDriftResult(false),
		makeDriftResult(false),
	}
	s := Summarize(results)
	if s.Total != 2 || s.InSync != 2 || s.Drifted != 0 || s.Missing != 0 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_WithDriftAndMissing(t *testing.T) {
	results := []DriftResult{
		makeDriftResult(false),
		makeDriftResult(false, FieldDiff{Field: "spec.replicas", Want: "3", Got: "1"}),
		makeDriftResult(true),
	}
	s := Summarize(results)
	if s.Total != 3 {
		t.Errorf("expected total=3, got %d", s.Total)
	}
	if s.InSync != 1 {
		t.Errorf("expected in-sync=1, got %d", s.InSync)
	}
	if s.Drifted != 1 {
		t.Errorf("expected drifted=1, got %d", s.Drifted)
	}
	if s.Missing != 1 {
		t.Errorf("expected missing=1, got %d", s.Missing)
	}
}

func TestSummary_HasDrift(t *testing.T) {
	clean := DriftSummary{Total: 2, InSync: 2}
	if clean.HasDrift() {
		t.Error("expected no drift")
	}
	dirty := DriftSummary{Total: 2, InSync: 1, Drifted: 1}
	if !dirty.HasDrift() {
		t.Error("expected drift")
	}
}

func TestSummary_String(t *testing.T) {
	s := DriftSummary{Total: 4, InSync: 2, Drifted: 1, Missing: 1}
	got := s.String()
	want := "total=4 in-sync=2 drifted=1 missing=1"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || s.HasDrift() {
		t.Errorf("expected empty summary, got %+v", s)
	}
}
