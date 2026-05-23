package differ

import (
	"bytes"
	"strings"
	"testing"
)

func TestReporter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)

	results := []DriftResult{
		{Key: "default/Deployment/nginx", Drifts: nil, Missing: false},
	}

	drifted := r.Report(results)
	if drifted {
		t.Error("expected no drift detected")
	}
	if !strings.Contains(buf.String(), "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", buf.String())
	}
}

func TestReporter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)

	results := []DriftResult{
		{
			Key:   "default/Deployment/nginx",
			Drifts: []string{"spec.replicas: want 3, got 1"},
		},
	}

	drifted := r.Report(results)
	if !drifted {
		t.Error("expected drift to be detected")
	}
	out := buf.String()
	if !strings.Contains(out, "[DRIFT]") {
		t.Errorf("expected [DRIFT] in output, got: %s", out)
	}
	if !strings.Contains(out, "spec.replicas") {
		t.Errorf("expected drift detail in output, got: %s", out)
	}
}

func TestReporter_MissingResource(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)

	results := []DriftResult{
		{Key: "staging/Service/api", Missing: true},
	}

	drifted := r.Report(results)
	if !drifted {
		t.Error("expected drift to be detected for missing resource")
	}
	if !strings.Contains(buf.String(), "[MISSING]") {
		t.Errorf("expected [MISSING] in output, got: %s", buf.String())
	}
}

func TestNewReporter_NilWriter(t *testing.T) {
	r := NewReporter(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
	// Should not panic when reporting
	results := []DriftResult{{Key: "default/ConfigMap/cfg"}}
	r.Report(results)
}
