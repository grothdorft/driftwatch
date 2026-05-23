package differ

import (
	"fmt"
	"io"
)

// Reporter formats and writes drift results to an output writer.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter that writes to the given writer.
func NewReporter(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Report writes a human-readable summary of drift results.
// Returns the number of drifted resources.
func (r *Reporter) Report(results []DriftResult) int {
	driftCount := 0
	for _, res := range results {
		if !res.Drifted {
			fmt.Fprintf(r.out, "[OK]      %s\n", res.Key)
			continue
		}
		driftCount++
		fmt.Fprintf(r.out, "[DRIFT]   %s\n", res.Key)
		for _, d := range res.Diffs {
			fmt.Fprintf(r.out, "            field: %s\n", d.Field)
			fmt.Fprintf(r.out, "              repo: %v\n", d.RepoVal)
			fmt.Fprintf(r.out, "              live: %v\n", d.LiveVal)
		}
	}
	fmt.Fprintf(r.out, "\nSummary: %d/%d resources drifted.\n", driftCount, len(results))
	return driftCount
}
