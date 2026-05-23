package differ

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// DriftResult holds the comparison result for a single manifest.
type DriftResult struct {
	Key    string
	Drifts []string
	Missing bool
}

// Reporter formats and writes drift results to an output stream.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// Report writes a human-readable drift summary for the given results.
// Returns true if any drift was detected.
func (r *Reporter) Report(results []DriftResult) bool {
	driftDetected := false
	for _, res := range results {
		if res.Missing {
			driftDetected = true
			fmt.Fprintf(r.out, "[MISSING] %s not found in live cluster\n", res.Key)
			continue
		}
		if len(res.Drifts) == 0 {
			fmt.Fprintf(r.out, "[OK]      %s\n", res.Key)
			continue
		}
		driftDetected = true
		fmt.Fprintf(r.out, "[DRIFT]   %s\n", res.Key)
		for _, d := range res.Drifts {
			fmt.Fprintf(r.out, "          - %s\n", strings.TrimSpace(d))
		}
	}
	return driftDetected
}

// Summary writes a single-line summary of the overall drift results,
// including counts of OK, drifted, and missing manifests.
func (r *Reporter) Summary(results []DriftResult) {
	ok, drifted, missing := 0, 0, 0
	for _, res := range results {
		switch {
		case res.Missing:
			missing++
		case len(res.Drifts) > 0:
			drifted++
		default:
			ok++
		}
	}
	fmt.Fprintf(r.out, "\nSummary: %d OK, %d drifted, %d missing (total: %d)\n",
		ok, drifted, missing, len(results))
}

// FilterDrifted returns only the results that have drift or are missing,
// making it easy to process or display only the problematic manifests.
func FilterDrifted(results []DriftResult) []DriftResult {
	var out []DriftResult
	for _, res := range results {
		if res.Missing || len(res.Drifts) > 0 {
			out = append(out, res)
		}
	}
	return out
}
