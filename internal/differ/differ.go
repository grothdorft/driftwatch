package differ

import (
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// DriftResult holds the comparison result for a single manifest.
type DriftResult struct {
	Key     string
	Drifted bool
	Diffs   []FieldDiff
}

// FieldDiff describes a single field-level difference.
type FieldDiff struct {
	Field    string
	RepoVal  interface{}
	LiveVal  interface{}
}

// Compare compares repo manifests against live manifests and returns drift results.
// repoManifests and liveManifests are keyed by ManifestKey (group/version/kind/namespace/name).
func Compare(
	repoManifests map[string]*unstructured.Unstructured,
	liveManifests map[string]*unstructured.Unstructured,
) []DriftResult {
	keys := make([]string, 0, len(repoManifests))
	for k := range repoManifests {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]DriftResult, 0, len(keys))
	for _, key := range keys {
		repo := repoManifests[key]
		live, exists := liveManifests[key]
		if !exists {
			results = append(results, DriftResult{
				Key:     key,
				Drifted: true,
				Diffs: []FieldDiff{{
					Field:   "<existence>",
					RepoVal: "present",
					LiveVal: "missing",
				}},
			})
			continue
		}
		diffs := diffObjects(repo.Object, live.Object, "")
		results = append(results, DriftResult{
			Key:     key,
			Drifted: len(diffs) > 0,
			Diffs:   diffs,
		})
	}
	return results
}

// diffObjects recursively compares two unstructured objects and returns field diffs.
func diffObjects(repo, live map[string]interface{}, prefix string) []FieldDiff {
	var diffs []FieldDiff
	for k, repoVal := range repo {
		field := fieldPath(prefix, k)
		liveVal, ok := live[k]
		if !ok {
			diffs = append(diffs, FieldDiff{Field: field, RepoVal: repoVal, LiveVal: nil})
			continue
		}
		repoMap, repoIsMap := repoVal.(map[string]interface{})
		liveMap, liveIsMap := liveVal.(map[string]interface{})
		if repoIsMap && liveIsMap {
			diffs = append(diffs, diffObjects(repoMap, liveMap, field)...)
		} else if fmt.Sprintf("%v", repoVal) != fmt.Sprintf("%v", liveVal) {
			diffs = append(diffs, FieldDiff{Field: field, RepoVal: repoVal, LiveVal: liveVal})
		}
	}
	return diffs
}

func fieldPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
