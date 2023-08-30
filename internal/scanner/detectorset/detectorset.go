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
	BuiltinAndSharedRuleIDs() []string
	TopLevelRuleIDs() []string
	DetectAt(
		node *tree.Node,
		ruleID string,
		detectorContext detectortypes.Context,
	) (*Result, error)
}

type detectorSet struct {
	builtinAndSharedRuleIDs,
	topLevelRuleIDs []string
	querySet  *query.Set
	detectors map[string]detectortypes.Detector
}

func New(
	querySet *query.Set,
	builtinDetectors []detectortypes.Detector,
	rules map[string]*settings.Rule,
	language language.Language,
) (Set, error) {
	relevantRules, presenceRules := getRelevantRules(rules, language.ID())
	builtinAndSharedRuleIDs, topLevelRuleIDs := findNotableRuleIDs(builtinDetectors, relevantRules, presenceRules)

	detectors, err := createDetectors(language, querySet, builtinDetectors, relevantRules)
	if err != nil {
		return nil, err
	}

	if err = querySet.Compile(); err != nil {
		return nil, fmt.Errorf("error compiling query set: %w", err)
	}

	return &detectorSet{
		builtinAndSharedRuleIDs: builtinAndSharedRuleIDs,
		topLevelRuleIDs:         topLevelRuleIDs,
		querySet:                querySet,
		detectors:               detectors,
	}, nil
}

func (set *detectorSet) BuiltinAndSharedRuleIDs() []string {
	return set.builtinAndSharedRuleIDs
}

func (set *detectorSet) TopLevelRuleIDs() []string {
	return set.topLevelRuleIDs
}

func (set *detectorSet) DetectAt(
	node *tree.Node,
	ruleID string,
	detectorContext detectortypes.Context,
) (*Result, error) {
	detector, err := set.lookupDetector(ruleID)
	if err != nil {
		return nil, err
	}

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
			RuleID:    ruleID,
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

	if ruleDetector.SanitizerRuleID() == "" {
		return false, nil
	}

	sanitizerDetections, err := detectorContext.Scan(
		node,
		ruleDetector.SanitizerRuleID(),
		settings.CURSOR_STRICT_SCOPE,
	)
	if err != nil {
		return false, err
	}

	return len(sanitizerDetections) != 0, nil
}

func (set *detectorSet) lookupDetector(ruleID string) (detectortypes.Detector, error) {
	detector, ok := set.detectors[ruleID]
	if !ok {
		return nil, fmt.Errorf("detector for rule '%s' not registered", ruleID)
	}

	return detector, nil
}

func getRelevantRules(
	rules map[string]*settings.Rule,
	languageID string,
) (map[string]*settings.Rule, set.Set[string]) {
	relevantRules := make(map[string]*settings.Rule)
	presenceRules := set.New[string]()

	for ruleID, rule := range rules {
		if !slices.Contains(rule.Languages, languageID) {
			continue
		}

		relevantRules[ruleID] = rule

		if rule.Trigger.RequiredDetection != nil {
			presenceRules.Add(*rule.Trigger.RequiredDetection)
		}
	}

	return relevantRules, presenceRules
}

func findNotableRuleIDs(
	builtinDetectors []detectortypes.Detector,
	relevantRules map[string]*settings.Rule,
	presenceRules set.Set[string],
) ([]string, []string) {
	var builtinAndSharedRuleIDs, topLevelRuleIDs []string

	for _, detector := range builtinDetectors {
		builtinAndSharedRuleIDs = append(builtinAndSharedRuleIDs, detector.RuleID())
	}

	for ruleID, rule := range relevantRules {
		if rule.Type == customdetectors.TypeShared {
			builtinAndSharedRuleIDs = append(builtinAndSharedRuleIDs, ruleID)
			continue
		}

		if !rule.IsAuxilary || presenceRules.Has(ruleID) {
			topLevelRuleIDs = append(topLevelRuleIDs, ruleID)
		}
	}

	return builtinAndSharedRuleIDs, topLevelRuleIDs
}

func createDetectors(
	language language.Language,
	querySet *query.Set,
	builtinDetectors []detectortypes.Detector,
	relevantRules map[string]*settings.Rule,
) (map[string]detectortypes.Detector, error) {
	detectors := make(map[string]detectortypes.Detector)

	for _, detector := range builtinDetectors {
		addDetector(detectors, detector)
	}

	for ruleID, rule := range relevantRules {
		detector, err := customrule.New(language, querySet, ruleID, rule.SanitizerRuleID, rule.Patterns)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s detector: %w", ruleID, err)
		}

		addDetector(detectors, detector)
	}

	return detectors, nil
}

func addDetector(detectorMap map[string]detectortypes.Detector, detector detectortypes.Detector) error {
	ruleID := detector.RuleID()

	if _, existing := detectorMap[ruleID]; existing {
		return fmt.Errorf("duplicate detector for rule '%s'", ruleID)
	}

	detectorMap[ruleID] = detector

	return nil
}
