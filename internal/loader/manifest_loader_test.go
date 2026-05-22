package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/loader"
)

const sampleManifest = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
spec:
  replicas: 3
`

func writeTempYAML(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := writeTempYAML(t, dir, "deploy.yaml", sampleManifest)

	m, err := loader.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Raw["kind"] != "Deployment" {
		t.Errorf("expected kind Deployment, got %v", m.Raw["kind"])
	}
	if m.Source != path {
		t.Errorf("expected source %s, got %s", path, m.Source)
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := loader.LoadFromFile("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromDir(t *testing.T) {
	dir := t.TempDir()
	writeTempYAML(t, dir, "deploy.yaml", sampleManifest)
	writeTempYAML(t, dir, "service.yml", `
apiVersion: v1
kind: Service
metadata:
  name: my-svc
  namespace: default
`)
	// non-yaml file should be ignored
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("ignore me"), 0644)

	manifests, err := loader.LoadFromDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests) != 2 {
		t.Errorf("expected 2 manifests, got %d", len(manifests))
	}
}

func TestManifestKey(t *testing.T) {
	dir := t.TempDir()
	path := writeTempYAML(t, dir, "deploy.yaml", sampleManifest)

	m, _ := loader.LoadFromFile(path)
	key, err := loader.ManifestKey(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "apps/v1/Deployment/default/my-app"
	if key != expected {
		t.Errorf("expected key %q, got %q", expected, key)
	}
}