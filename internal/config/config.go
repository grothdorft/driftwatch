package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the driftwatch runtime configuration.
type Config struct {
	// ManifestDir is the path to the directory containing source-of-truth YAML files.
	ManifestDir string `yaml:"manifestDir"`

	// Namespace restricts drift detection to a specific Kubernetes namespace.
	// If empty, all namespaces are checked.
	Namespace string `yaml:"namespace"`

	// Kubeconfig is the path to the kubeconfig file.
	// Defaults to in-cluster config when empty.
	Kubeconfig string `yaml:"kubeconfig"`

	// FailOnDrift causes the process to exit with a non-zero status when drift is detected.
	FailOnDrift bool `yaml:"failOnDrift"`

	// OutputFormat controls report output. Valid values: "text", "json".
	OutputFormat string `yaml:"outputFormat"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ManifestDir:  "./manifests",
		Namespace:    "",
		Kubeconfig:   "",
		FailOnDrift:  false,
		OutputFormat: "text",
	}
}

// LoadFromFile reads a YAML config file from the given path and returns a Config.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the Config contains valid values.
func (c *Config) Validate() error {
	if c.ManifestDir == "" {
		return fmt.Errorf("manifestDir must not be empty")
	}
	switch c.OutputFormat {
	case "text", "json":
		// valid
	default:
		return fmt.Errorf("unsupported outputFormat %q: must be \"text\" or \"json\"", c.OutputFormat)
	}
	return nil
}
