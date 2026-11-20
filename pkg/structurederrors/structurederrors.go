// Package structurederrors records problems with rules in a structured form, so that tools (such as
// an agent writing rules) can act on them without parsing log output.
//
// Recording is opt-in: a nil *Collector is valid and ignores everything, so call sites do not need to
// check whether structured errors are enabled. Recording never changes how the scan behaves.
package structurederrors

import (
	"cmp"
	"slices"
	"sync"
)

// SchemaVersion is bumped on any breaking change to the Report structure
const SchemaVersion = "1"

type Severity string

const (
	// SeverityError means the rule was not run. When the error stopped the scan, the report has
	// no findings.
	SeverityError Severity = "error"
	// SeverityWarning means the rule is valid but something about it needs attention
	SeverityWarning Severity = "warning"
)

// Code identifies the kind of problem. Codes are stable and intended to be matched on by tools;
// messages are for humans and may change.
type Code string

const (
	CodeRuleReadFailed         Code = "rule_read_failed"
	CodeRuleYAMLInvalid        Code = "rule_yaml_invalid"
	CodeRuleMissingMetadata    Code = "rule_missing_metadata"
	CodeRuleMissingID          Code = "rule_missing_id"
	CodeRuleDuplicateID        Code = "rule_duplicate_id"
	CodeRuleLanguageNotFound   Code = "rule_language_not_in_target"
	CodeImportUnknownRule      Code = "import_unknown_rule"
	CodeImportNotShared        Code = "import_not_shared"
	CodeFilterUnknownRule      Code = "filter_unknown_rule"
	CodeSanitizerUnknownRule   Code = "sanitizer_unknown_rule"
	CodeSharedRuleInvalidField Code = "shared_rule_invalid_field"
	CodePatternInvalid         Code = "pattern_invalid"
	CodeFilterInvalid          Code = "filter_invalid"
	// CodeRuleSetInvalid is a problem with the rules for a language that is not specific to one
	// pattern
	CodeRuleSetInvalid Code = "rule_set_invalid"
)

type Item struct {
	Severity Severity `json:"severity" yaml:"severity"`
	Code     Code     `json:"code" yaml:"code"`
	Message  string   `json:"message" yaml:"message"`
	// RuleID is the rule the problem belongs to. For an auxiliary rule this is the auxiliary ID, and
	// ParentRuleID is the rule that defines it.
	RuleID       string `json:"rule_id,omitempty" yaml:"rule_id,omitempty"`
	ParentRuleID string `json:"parent_rule_id,omitempty" yaml:"parent_rule_id,omitempty"`
	// RuleFile is the path of the rule's YAML file
	RuleFile string `json:"rule_file,omitempty" yaml:"rule_file,omitempty"`
	// PatternIndex is the 0-based index into the rule's `patterns` list
	PatternIndex *int   `json:"pattern_index,omitempty" yaml:"pattern_index,omitempty"`
	Pattern      string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Language     string `json:"language,omitempty" yaml:"language,omitempty"`
	// RelatedRuleID is another rule involved in the problem (eg. the missing import)
	RelatedRuleID string `json:"related_rule_id,omitempty" yaml:"related_rule_id,omitempty"`
}

// Report is the structured_errors section of a report
type Report struct {
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`
	Items         []Item `json:"items" yaml:"items"`
}

type ruleSource struct {
	file     string
	parentID string
}

// Collector gathers structured errors. It is safe for concurrent use. A nil Collector discards
// everything.
type Collector struct {
	mutex       sync.Mutex
	items       []Item
	ruleSources map[string]ruleSource
}

func New() *Collector {
	return &Collector{ruleSources: make(map[string]ruleSource)}
}

func (collector *Collector) Enabled() bool {
	return collector != nil
}

// RegisterRule records where a rule came from, so later items for it are filled in with its file
// (and parent, for auxiliary rules)
func (collector *Collector) RegisterRule(id, file, parentID string) {
	if collector == nil {
		return
	}

	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	collector.ruleSources[id] = ruleSource{file: file, parentID: parentID}
}

func (collector *Collector) Add(item Item) {
	if collector == nil {
		return
	}

	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	if source, exists := collector.ruleSources[item.RuleID]; exists {
		if item.RuleFile == "" {
			item.RuleFile = source.file
		}
		if item.ParentRuleID == "" {
			item.ParentRuleID = source.parentID
		}
	}

	collector.items = append(collector.items, item)
}

// Report returns the collected items in a stable order, or nil when the collector is disabled
func (collector *Collector) Report() *Report {
	if collector == nil {
		return nil
	}

	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	items := slices.Clone(collector.items)
	if items == nil {
		items = []Item{}
	}

	slices.SortStableFunc(items, func(a, b Item) int {
		return cmp.Or(
			cmp.Compare(a.RuleFile, b.RuleFile),
			cmp.Compare(a.RuleID, b.RuleID),
			cmp.Compare(patternIndex(a), patternIndex(b)),
			cmp.Compare(a.Language, b.Language),
		)
	})

	return &Report{SchemaVersion: SchemaVersion, Items: items}
}

func patternIndex(item Item) int {
	if item.PatternIndex == nil {
		return -1
	}

	return *item.PatternIndex
}

// Index returns a pointer to the given int, for use as a PatternIndex
func Index(index int) *int {
	return &index
}
