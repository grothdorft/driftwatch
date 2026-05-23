package differ

import "fmt"

// DriftSummary holds aggregate statistics from a drift comparison run.
type DriftSummary struct {
	Total    int
	Drifted  int
	Missing  int
	InSync   int
}

// Summarize iterates over a slice of DriftResults and returns a DriftSummary.
func Summarize(results []DriftResult) DriftSummary {
	s := DriftSummary{Total: len(results)}
	for _, r := range results {
		switch {
		case r.Missing:
			s.Missing++
		case len(r.Diffs) > 0:
			s.Drifted++
		default:
			s.InSync++
		}
	}
	return s
}

// String returns a human-readable one-line summary.
func (s DriftSummary) String() string {
	return fmt.Sprintf(
		"total=%d in-sync=%d drifted=%d missing=%d",
		s.Total, s.InSync, s.Drifted, s.Missing,
	)
}

// HasDrift returns true when any resource is drifted or missing.
func (s DriftSummary) HasDrift() bool {
	return s.Drifted > 0 || s.Missing > 0
}
