package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.ManifestDir != "./manifests" {
		t.Errorf("expected default manifestDir './manifests', got %q", cfg.ManifestDir)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default outputFormat 'text', got %q", cfg.OutputFormat)
	}
	if cfg.FailOnDrift {
		t.Error("expected default failOnDrift to be false")
	}
}

func TestLoadFromFile_Valid(t *testing.T) {
	path := writeTempConfig(t, `
manifestDir: /repo/k8s
namespace: production
failOnDrift: true
outputFormat: json
`)
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ManifestDir != "/repo/k8s" {
		t.Errorf("manifestDir: got %q, want '/repo/k8s'", cfg.ManifestDir)
	}
	if cfg.Namespace != "production" {
		t.Errorf("namespace: got %q, want 'production'", cfg.Namespace)
	}
	if !cfg.FailOnDrift {
		t.Error("expected failOnDrift true")
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("outputFormat: got %q, want 'json'", cfg.OutputFormat)
	}
}

func TestLoadFromFile_NotFound(t *testing.T) {
	_, err := config.LoadFromFile("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFromFile_InvalidFormat(t *testing.T) {
	path := writeTempConfig(t, `outputFormat: xml\n`)
	_, err := config.LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for unsupported outputFormat")
	}
}

func TestValidate_EmptyManifestDir(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.ManifestDir = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty manifestDir")
	}
}

func TestValidate_ValidFormats(t *testing.T) {
	for _, f := range []string{"text", "json"} {
		cfg := config.DefaultConfig()
		cfg.OutputFormat = f
		if err := cfg.Validate(); err != nil {
			t.Errorf("format %q should be valid, got: %v", f, err)
		}
	}
}
