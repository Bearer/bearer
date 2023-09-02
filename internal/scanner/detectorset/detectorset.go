package detectorset

import (
	"fmt"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectors/customrule"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

const ()

type Result struct {
	Detections []*detectortypes.Detection
	Sanitized  bool
}

type Set interface {
	DetectAt(
		node *tree.Node,
		rule *ruleset.Rule,
		detectorContext detectortypes.Context,
	) (*Result, error)
}

type detectorSet struct {
	detectors []detectortypes.Detector
}

func New(
	schemaClassifier *schema.Classifier,
	language language.Language,
	ruleSet *ruleset.Set,
	querySet *query.Set,
) (Set, error) {
	detectors := make([]detectortypes.Detector, len(ruleSet.Rules()))

	for _, detector := range language.NewBuiltInDetectors(schemaClassifier, querySet) {
		detectors[detector.Rule().Index()] = detector
	}

	for _, rule := range ruleSet.Rules() {
		if rule.Type() == ruleset.RuleTypeBuiltin {
			continue
		}

		detector, err := customrule.New(language, ruleSet, querySet, rule)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s detector: %w", rule.ID(), err)
		}

		detectors[rule.Index()] = detector
	}

	return &detectorSet{
		detectors: detectors,
	}, nil
}

func (set *detectorSet) DetectAt(
	node *tree.Node,
	rule *ruleset.Rule,
	detectorContext detectortypes.Context,
) (*Result, error) {
	detector := set.detectors[rule.Index()]

	if isSanitized, err := set.isSanitized(rule, node, detectorContext); isSanitized || err != nil {
		return &Result{Sanitized: true}, err
	}

	detectionsData, err := detector.DetectAt(node, detectorContext)
	if err != nil {
		return nil, err
	}

	detections := make([]*detectortypes.Detection, len(detectionsData))
	for i, data := range detectionsData {
		detections[i] = &detectortypes.Detection{
			RuleID:    rule.ID(),
			MatchNode: node,
			Data:      data,
		}
	}

	return &Result{Detections: detections}, nil
}

func (set *detectorSet) isSanitized(
	rule *ruleset.Rule,
	node *tree.Node,
	detectorContext detectortypes.Context,
) (bool, error) {
	sanitizerRule := rule.SanitizerRule()
	if sanitizerRule == nil {
		return false, nil
	}

	detections, err := detectorContext.Scan(node, sanitizerRule, settings.CURSOR_STRICT_SCOPE)
	if err != nil {
		return false, err
	}

	return len(detections) != 0, nil
}
