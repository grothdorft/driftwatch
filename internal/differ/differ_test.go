package differ_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/differ"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func makeUnstructured(obj map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: obj}
}

func TestCompare_NoDrift(t *testing.T) {
	obj := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": "cfg", "namespace": "default"},
		"data":       map[string]interface{}{"key": "value"},
	}
	repo := map[string]*unstructured.Unstructured{"v1/ConfigMap/default/cfg": makeUnstructured(obj)}
	live := map[string]*unstructured.Unstructured{"v1/ConfigMap/default/cfg": makeUnstructured(obj)}

	results := differ.Compare(repo, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got diffs: %+v", results[0].Diffs)
	}
}

func TestCompare_FieldDrift(t *testing.T) {
	repoObj := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": "cfg", "namespace": "default"},
		"data":       map[string]interface{}{"key": "repo-value"},
	}
	liveObj := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": "cfg", "namespace": "default"},
		"data":       map[string]interface{}{"key": "live-value"},
	}
	repo := map[string]*unstructured.Unstructured{"v1/ConfigMap/default/cfg": makeUnstructured(repoObj)}
	live := map[string]*unstructured.Unstructured{"v1/ConfigMap/default/cfg": makeUnstructured(liveObj)}

	results := differ.Compare(repo, live)
	if !results[0].Drifted {
		t.Fatal("expected drift but got none")
	}
	if len(results[0].Diffs) != 1 || results[0].Diffs[0].Field != "data.key" {
		t.Errorf("unexpected diffs: %+v", results[0].Diffs)
	}
}

func TestCompare_MissingInLive(t *testing.T) {
	repoObj := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": "missing", "namespace": "default"},
	}
	repo := map[string]*unstructured.Unstructured{"v1/ConfigMap/default/missing": makeUnstructured(repoObj)}
	live := map[string]*unstructured.Unstructured{}

	results := differ.Compare(repo, live)
	if !results[0].Drifted {
		t.Fatal("expected drift for missing live resource")
	}
	if results[0].Diffs[0].Field != "<existence>" {
		t.Errorf("expected existence diff, got: %+v", results[0].Diffs)
	}
}
