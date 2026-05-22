# driftwatch

Detects configuration drift between deployed Kubernetes manifests and their source-of-truth YAML files in a repo.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at a directory of manifest files and it will compare them against what is currently deployed in your cluster.

```bash
# Compare manifests in a local directory against the live cluster
driftwatch --dir ./manifests --namespace production

# Use a specific kubeconfig
driftwatch --dir ./manifests --namespace staging --kubeconfig ~/.kube/config

# Output results as JSON
driftwatch --dir ./manifests --namespace production --output json
```

**Example output:**

```
[DRIFT]  Deployment/api-server  image: repo/api:v1.2.0 → repo/api:v1.3.1
[DRIFT]  ConfigMap/app-config   replicas: 3 → 5
[OK]     Service/api-server
```

Exit code `1` is returned if any drift is detected, making it easy to integrate into CI pipelines.

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--dir` | `.` | Path to directory containing YAML manifests |
| `--namespace` | `default` | Kubernetes namespace to compare against |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |
| `--output` | `text` | Output format: `text` or `json` |

---

## License

MIT © 2024 yourusername