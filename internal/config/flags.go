package config

import (
	"flag"
	"fmt"
)

// Flags holds CLI flag values that can override file-based config.
type Flags struct {
	ConfigFile   string
	ManifestDir  string
	Namespace    string
	Kubeconfig   string
	FailOnDrift  bool
	OutputFormat string
}

// RegisterFlags binds Flags fields to the provided FlagSet.
func RegisterFlags(fs *flag.FlagSet, f *Flags) {
	fs.StringVar(&f.ConfigFile, "config", "", "Path to driftwatch config file (YAML)")
	fs.StringVar(&f.ManifestDir, "manifest-dir", "", "Directory containing source-of-truth manifests")
	fs.StringVar(&f.Namespace, "namespace", "", "Kubernetes namespace to inspect (empty = all)")
	fs.StringVar(&f.Kubeconfig, "kubeconfig", "", "Path to kubeconfig (defaults to in-cluster)")
	fs.BoolVar(&f.FailOnDrift, "fail-on-drift", false, "Exit non-zero when drift is detected")
	fs.StringVar(&f.OutputFormat, "output", "", "Output format: text or json")
}

// Merge applies non-zero flag values on top of cfg, returning the merged result.
func Merge(cfg *Config, f *Flags) (*Config, error) {
	out := *cfg // shallow copy

	if f.ManifestDir != "" {
		out.ManifestDir = f.ManifestDir
	}
	if f.Namespace != "" {
		out.Namespace = f.Namespace
	}
	if f.Kubeconfig != "" {
		out.Kubeconfig = f.Kubeconfig
	}
	if f.FailOnDrift {
		out.FailOnDrift = true
	}
	if f.OutputFormat != "" {
		out.OutputFormat = f.OutputFormat
	}

	if err := out.Validate(); err != nil {
		return nil, fmt.Errorf("merged config invalid: %w", err)
	}
	return &out, nil
}
