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
