package stats

import (
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Stats struct {
	timings map[string]time.Duration
}

func New() *Stats {
	return &Stats{timings: make(map[string]time.Duration)}
}

func (stats *Stats) Record(ruleID string, startTime time.Time) {
	duration := time.Since(startTime)
	stats.timings[ruleID] += duration
}

func (stats *Stats) String() string {
	sortedRuleIDs := maps.Keys(stats.timings)
	slices.SortFunc(sortedRuleIDs, func(a, b string) bool {
		return stats.timings[a] > stats.timings[b]
	})

	var s strings.Builder
	for i, ruleID := range sortedRuleIDs {
		if i != 0 {
			s.WriteString("\n")
		}

		s.WriteString(ruleID)
		s.WriteString(": ")
		s.WriteString(stats.timings[ruleID].String())
	}

	return s.String()
}
