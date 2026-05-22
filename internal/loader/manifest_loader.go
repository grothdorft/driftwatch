// Package loader provides utilities for loading Kubernetes manifests
// from local YAML files and from a live cluster.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// Manifest represents a raw Kubernetes manifest decoded from YAML.
type Manifest struct {
	// Source is the file path or cluster identifier where this manifest was loaded from.
	Source string
	// Raw holds the decoded manifest as a generic map.
	Raw map[string]interface{}
}

// LoadFromFile reads a single YAML file and returns a Manifest.
func LoadFromFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unmarshalling yaml from %s: %w", path, err)
	}

	return &Manifest{Source: path, Raw: raw}, nil
}

// LoadFromDir walks a directory and loads all *.yaml and *.yml files.
func LoadFromDir(dir string) ([]*Manifest, error) {
	var manifests []*Manifest

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		m, err := LoadFromFile(path)
		if err != nil {
			return fmt.Errorf("loading manifest %s: %w", path, err)
		}
		manifests = append(manifests, m)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", dir, err)
	}

	return manifests, nil
}

// ManifestKey returns a unique identifier for a manifest based on
// its apiVersion, kind, namespace, and name.
func ManifestKey(m *Manifest) (string, error) {
	apiVersion, _ := m.Raw["apiVersion"].(string)
	kind, _ := m.Raw["kind"].(string)
	if apiVersion == "" || kind == "" {
		return "", fmt.Errorf("manifest missing apiVersion or kind in %s", m.Source)
	}

	meta, _ := m.Raw["metadata"].(map[string]interface{})
	name, _ := meta["name"].(string)
	namespace, _ := meta["namespace"].(string)

	return fmt.Sprintf("%s/%s/%s/%s", apiVersion, kind, namespace, name), nil
}