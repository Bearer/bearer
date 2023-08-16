package set

import (
	"fmt"
	"slices"

	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/detector/implementation/custom"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/util/set"
)

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
	languageName string,
	lang types.Language,
) (detectortypes.DetectorSet, error) {
	relevantRules, presenceRules := getRelevantRules(rules, languageName)
	builtinAndSharedRuleIDs, topLevelRuleIDs := findNotableRuleIDs(builtinDetectors, relevantRules, presenceRules)

	detectors, err := createDetectors(lang, querySet, builtinDetectors, relevantRules)
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

func (set *detectorSet) NestedDetections(detectorType string) (bool, error) {
	detector, err := set.lookupDetector(detectorType)
	if err != nil {
		return false, err
	}

	return detector.NestedDetections(), nil
}

func (set *detectorSet) DetectAt(
	node *tree.Node,
	detectorType string,
	evaluationState detectortypes.EvaluationState,
) ([]*detection.Detection, error) {
	detector, err := set.lookupDetector(detectorType)
	if err != nil {
		return nil, err
	}

	detectionsData, err := detector.DetectAt(node, evaluationState)
	if err != nil {
		return nil, err
	}

	detections := make([]*detection.Detection, len(detectionsData))
	for i, data := range detectionsData {
		detections[i] = &detection.Detection{
			RuleID:    detectorType,
			MatchNode: node,
			Data:      data,
		}
	}

	return detections, nil
}

func (set *detectorSet) lookupDetector(detectorType string) (detectortypes.Detector, error) {
	detector, ok := set.detectors[detectorType]
	if !ok {
		return nil, fmt.Errorf("detector type '%s' not registered", detectorType)
	}

	return detector, nil
}

func getRelevantRules(
	rules map[string]*settings.Rule,
	languageName string,
) (map[string]*settings.Rule, set.Set[string]) {
	relevantRules := make(map[string]*settings.Rule)
	presenceRules := set.New[string]()

	for ruleID, rule := range rules {
		if !slices.Contains(rule.Languages, languageName) {
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
		builtinAndSharedRuleIDs = append(builtinAndSharedRuleIDs, detector.Name())
	}

	for ruleID, rule := range relevantRules {
		if rule.Type == customdetectors.TypeShared {
			builtinAndSharedRuleIDs = append(builtinAndSharedRuleIDs, ruleID)
		}

		if !rule.IsAuxilary || presenceRules.Has(ruleID) {
			topLevelRuleIDs = append(topLevelRuleIDs, ruleID)
		}
	}

	return builtinAndSharedRuleIDs, topLevelRuleIDs
}

func createDetectors(
	lang types.Language,
	querySet *query.Set,
	builtinDetectors []detectortypes.Detector,
	relevantRules map[string]*settings.Rule,
) (map[string]detectortypes.Detector, error) {
	detectors := make(map[string]detectortypes.Detector)

	for _, detector := range builtinDetectors {
		addDetector(detectors, detector)
	}

	for ruleID, rule := range relevantRules {
		detector, err := custom.New(
			lang,
			querySet,
			ruleID,
			rule.Patterns,
			relevantRules,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s detector: %w", ruleID, err)
		}

		addDetector(detectors, detector)
	}

	return detectors, nil
}

func addDetector(detectorMap map[string]detectortypes.Detector, detector detectortypes.Detector) error {
	name := detector.Name()

	if _, existing := detectorMap[name]; existing {
		return fmt.Errorf("duplicate detector '%s'", name)
	}

	detectorMap[name] = detector

	return nil
}
