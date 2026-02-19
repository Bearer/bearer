package stats

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

const maxSlowFiles = 10
const maxSlowRules = 10

type FailedReason int

const (
	TimeoutReason FailedReason = iota
	KilledAfterTimeoutReason
	MemoryLimitReason
	UnexpectedReason
)

type slowFile struct {
	filename string
	duration time.Duration
}

type failedFile struct {
	filename    string
	reason      FailedReason
	duration    time.Duration
	memoryUsage uint64
}

type FileStats struct {
	rules map[string]time.Duration
}

type fileStatsJSON struct {
	Rules map[string]time.Duration
}

type Stats struct {
	rules             map[string]time.Duration
	slowFiles         []slowFile
	totalFileDuration time.Duration
	failedFiles       []failedFile
	fileMutex         sync.Mutex
}

func NewFileStats() *FileStats {
	return &FileStats{rules: make(map[string]time.Duration)}
}

func (stats *FileStats) Rule(ruleID string, startTime time.Time) {
	if stats == nil {
		return
	}

	duration := time.Since(startTime)
	stats.rules[ruleID] += duration
}

func (stats *FileStats) MarshalJSON() ([]byte, error) {
	var statsJSON *fileStatsJSON
	if stats != nil {
		statsJSON = &fileStatsJSON{Rules: stats.rules}
	}

	return json.Marshal(statsJSON)
}

func (stats *FileStats) UnmarshalJSON(input []byte) error {
	var statsJSON *fileStatsJSON

	if err := json.Unmarshal(input, &statsJSON); err != nil {
		return err
	}

	if statsJSON != nil {
		*stats = FileStats{rules: statsJSON.Rules}
	}

	return nil
}

func New() *Stats {
	return &Stats{rules: make(map[string]time.Duration)}
}

func (stats *Stats) File(filename string, startTime time.Time) time.Duration {
	duration := time.Since(startTime)

	if stats == nil {
		return duration
	}

	stats.totalFileDuration += duration

	if len(stats.slowFiles) < maxSlowFiles {
		stats.slowFiles = append(stats.slowFiles, slowFile{
			filename: filename,
			duration: duration,
		})
	}

	fastestFileIndex := 0
	for i, file := range stats.slowFiles {
		if file.duration < stats.slowFiles[fastestFileIndex].duration {
			fastestFileIndex = i
		}
	}

	if duration > stats.slowFiles[fastestFileIndex].duration {
		stats.slowFiles[fastestFileIndex].filename = filename
		stats.slowFiles[fastestFileIndex].duration = duration
	}

	return duration
}

func (stats *Stats) FileFailed(filename string, reason FailedReason, startTime time.Time, memoryUsage uint64) {
	if stats == nil {
		return
	}

	stats.failedFiles = append(stats.failedFiles, failedFile{
		filename:    filename,
		reason:      reason,
		duration:    time.Since(startTime),
		memoryUsage: memoryUsage,
	})
}

func (stats *Stats) AddFileStats(fileStats *FileStats) {
	if stats == nil {
		return
	}

	stats.fileMutex.Lock()
	defer stats.fileMutex.Unlock()

	for ruleID, duration := range fileStats.rules {
		stats.rules[ruleID] += duration
	}
}

func (stats *Stats) String() string {
	var s strings.Builder

	stats.reportSlowestFiles(&s)
	stats.reportSlowestRules(&s)
	stats.reportFailedFiles(&s)

	return s.String()
}

func (stats *Stats) reportSlowestFiles(writer io.StringWriter) {
	writer.WriteString(fmt.Sprintf( //nolint:errcheck
		"Slowest files (total runtime %s):\n",
		stats.totalFileDuration.Truncate(time.Millisecond)),
	)
	slices.SortFunc(stats.slowFiles, func(a, b slowFile) int {
		if a.duration == b.duration {
			return strings.Compare(a.filename, b.filename)
		}

		return int(b.duration - a.duration)
	})

	for _, file := range stats.slowFiles {
		percentage := (float64(file.duration) / float64(stats.totalFileDuration)) * 100
		writer.WriteString(fmt.Sprintf( //nolint:errcheck
			"  - %s [%s %.2f%%]\n",
			file.filename,
			file.duration.Truncate(time.Millisecond),
			percentage,
		))
	}
}

func (stats *Stats) reportSlowestRules(writer io.StringWriter) {
	var totalRuleDuration time.Duration
	for _, ruleDuration := range stats.rules {
		totalRuleDuration += ruleDuration
	}

	writer.WriteString(fmt.Sprintf( //nolint:errcheck
		"\nSlowest rules (total runtime %s):\n",
		totalRuleDuration.Truncate(time.Millisecond),
	))
	sortedRuleIDs := slices.Collect(maps.Keys(stats.rules))
	slices.SortFunc(sortedRuleIDs, func(a, b string) int {
		return int(stats.rules[b] - stats.rules[a])
	})

	numSlowRules := maxSlowRules
	if numSlowRules > len(sortedRuleIDs) {
		numSlowRules = len(sortedRuleIDs)
	}

	for _, ruleID := range sortedRuleIDs[:numSlowRules] {
		ruleDuration := stats.rules[ruleID]
		percentage := (float64(ruleDuration) / float64(totalRuleDuration)) * 100
		writer.WriteString(fmt.Sprintf( //nolint:errcheck
			"  - %s [%s %.2f%%]\n",
			ruleID,
			ruleDuration.Truncate(time.Millisecond),
			percentage,
		))
	}
}

func (stats *Stats) reportFailedFiles(writer io.StringWriter) {
	if len(stats.failedFiles) == 0 {
		return
	}

	writer.WriteString("\nFailed files:\n") //nolint:errcheck
	slices.SortFunc(stats.failedFiles, func(a, b failedFile) int {
		return strings.Compare(a.filename, b.filename)
	})

	for _, file := range stats.failedFiles {
		if file.reason == MemoryLimitReason {
			writer.WriteString(fmt.Sprintf( //nolint:errcheck
				"  - %s [%s %s]\n",
				file.filename,
				file.reason,
				humanize.Bytes(file.memoryUsage),
			))
			continue
		}

		writer.WriteString(fmt.Sprintf( //nolint:errcheck
			"  - %s [%s %s]\n",
			file.filename,
			file.reason,
			file.duration.Truncate(time.Millisecond),
		))
	}
}

func (reason FailedReason) String() string {
	switch reason {
	case TimeoutReason:
		return "Time limit exceeded"
	case KilledAfterTimeoutReason:
		return "Time limit exceeded (killed)"
	case MemoryLimitReason:
		return "Memory limit exceeded"
	case UnexpectedReason:
		return "Unexpected (see debug logs)"
	default:
		panic(fmt.Sprintf("unknown stats file failed reason %d", reason))
	}
}
