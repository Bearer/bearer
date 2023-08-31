package detectorset

import (
	"fmt"
	"slices"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/report/customdetectors"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/util/set"
)

type Result struct {
	Detections []*detectortypes.Detection
	Sanitized  bool
}

type Set interface {
	BuiltinAndSharedDetectorIDs() []int
	TopLevelDetectorIDs() []int
	RuleIDFor(detectorID int) string
	DetectorIDFor(ruleID string) int
	DetectAt(
		node *tree.Node,
		detectorID int,
		detectorContext detectortypes.Context,
	) (*Result, error)
}

type detectorSet struct {
	// FIXME: remove this
	detectorIDByRuleID map[string]int
	builtinAndSharedDetectorIDs,
	topLevelDetectorIDs []int
	querySet  *query.Set
	detectors []detectortypes.Detector
}

func New(
	querySet *query.Set,
	builtinDetectors []detectortypes.Detector,
	rules map[string]*settings.Rule,
	language language.Language,
) (Set, error) {
	relevantRules, presenceRules := getRelevantRules(rules, language.ID())
	detectors, detectorIDByRuleID, err := createDetectors(language, querySet, builtinDetectors, relevantRules)
	if err != nil {
		return nil, err
	}

	builtinAndSharedDetectorIDs, topLevelDetectorIDs := findNotableIDs(
		detectorIDByRuleID,
		builtinDetectors,
		relevantRules,
		presenceRules,
	)

	if err = querySet.Compile(); err != nil {
		return nil, fmt.Errorf("error compiling query set: %w", err)
	}

	return &detectorSet{
		detectorIDByRuleID:          detectorIDByRuleID,
		builtinAndSharedDetectorIDs: builtinAndSharedDetectorIDs,
		topLevelDetectorIDs:         topLevelDetectorIDs,
		querySet:                    querySet,
		detectors:                   detectors,
	}, nil
}

func (set *detectorSet) BuiltinAndSharedDetectorIDs() []int {
	return set.builtinAndSharedDetectorIDs
}

func (set *detectorSet) TopLevelDetectorIDs() []int {
	return set.topLevelDetectorIDs
}

func (set *detectorSet) RuleIDFor(detectorID int) string {
	return set.detectors[detectorID].RuleID()
}

func (set *detectorSet) DetectorIDFor(ruleID string) int {
	detectorID, exists := set.detectorIDByRuleID[ruleID]
	if !exists {
		panic(fmt.Sprintf("unknown rule %s", ruleID))
	}

	return detectorID
}

func (set *detectorSet) DetectAt(
	node *tree.Node,
	detectorID int,
	detectorContext detectortypes.Context,
) (*Result, error) {
	detector := set.detectors[detectorID]

	if isSanitized, err := set.isSanitized(detector, node, detectorContext); isSanitized || err != nil {
		return &Result{Sanitized: true}, err
	}

	detectionsData, err := detector.DetectAt(node, detectorContext)
	if err != nil {
		return nil, err
	}

	detections := make([]*detectortypes.Detection, len(detectionsData))
	for i, data := range detectionsData {
		detections[i] = &detectortypes.Detection{
			RuleID:    detector.RuleID(),
			MatchNode: node,
			Data:      data,
		}
	}

	return &Result{Detections: detections}, nil
}

func (set *detectorSet) isSanitized(
	detector detectortypes.Detector,
	node *tree.Node,
	detectorContext detectortypes.Context,
) (bool, error) {
	ruleDetector, isCustomRule := detector.(*customrule.Detector)
	if !isCustomRule {
		return false, nil
	}

	if ruleDetector.SanitizerDetectorID() == -1 {
		return false, nil
	}

	sanitizerDetections, err := detectorContext.Scan(
		node,
		ruleDetector.SanitizerDetectorID(),
		settings.CURSOR_STRICT_SCOPE,
	)
	if err != nil {
		return false, err
	}

	return len(sanitizerDetections) != 0, nil
}

func getRelevantRules(
	rules map[string]*settings.Rule,
	languageID string,
) ([]*settings.Rule, set.Set[string]) {
	var relevantRules []*settings.Rule
	presenceRules := set.New[string]()

	for _, rule := range rules {
		if !slices.Contains(rule.Languages, languageID) {
			continue
		}

		relevantRules = append(relevantRules, rule)

		if rule.Trigger.RequiredDetection != nil {
			presenceRules.Add(*rule.Trigger.RequiredDetection)
		}
	}

	return relevantRules, presenceRules
}

func findNotableIDs(
	detectorIDByRuleID map[string]int,
	builtinDetectors []detectortypes.Detector,
	relevantRules []*settings.Rule,
	presenceRules set.Set[string],
) ([]int, []int) {
	var builtinAndSharedDetectorIDs, topLevelDetectorIDs []int

	for _, detector := range builtinDetectors {
		builtinAndSharedDetectorIDs = append(builtinAndSharedDetectorIDs, detectorIDByRuleID[detector.RuleID()])
	}

	for _, rule := range relevantRules {
		if rule.Type == customdetectors.TypeShared {
			builtinAndSharedDetectorIDs = append(builtinAndSharedDetectorIDs, detectorIDByRuleID[rule.Id])
			continue
		}

		if !rule.IsAuxilary || presenceRules.Has(rule.Id) {
			topLevelDetectorIDs = append(topLevelDetectorIDs, detectorIDByRuleID[rule.Id])
		}
	}

	return builtinAndSharedDetectorIDs, topLevelDetectorIDs
}

func createDetectors(
	language language.Language,
	querySet *query.Set,
	builtinDetectors []detectortypes.Detector,
	relevantRules []*settings.Rule,
) ([]detectortypes.Detector, map[string]int, error) {
	detectorIDByRuleID, err := allocateIDs(builtinDetectors, relevantRules)
	if err != nil {
		return nil, nil, err
	}

	detectors := builtinDetectors
	for _, rule := range relevantRules {
		detector, err := customrule.New(
			language,
			querySet,
			detectorIDByRuleID,
			rule.Id,
			rule.SanitizerRuleID,
			rule.Patterns,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create %s detector: %w", rule.Id, err)
		}

		detectors = append(detectors, detector)
	}

	return detectors, detectorIDByRuleID, nil
}

func allocateIDs(
	builtinDetectors []detectortypes.Detector,
	relevantRules []*settings.Rule,
) (map[string]int, error) {
	result := make(map[string]int)

	for i, detector := range builtinDetectors {
		ruleID := detector.RuleID()
		if _, existing := result[ruleID]; existing {
			return nil, fmt.Errorf("duplicate built-in detector for rule '%s'", ruleID)
		}

		result[ruleID] = i
	}

	for i, rule := range relevantRules {
		ruleID := rule.Id
		if _, existing := result[ruleID]; existing {
			return nil, fmt.Errorf("duplicate detector for rule '%s'", ruleID)
		}

		result[ruleID] = i + len(builtinDetectors)
	}

	return result, nil
}
