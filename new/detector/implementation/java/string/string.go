package string

import (
	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"

	generictypes "github.com/bearer/bearer/new/detector/implementation/generic/types"
	languagetypes "github.com/bearer/bearer/new/language/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(lang languagetypes.Language) (types.Detector, error) {
	return &stringDetector{}, nil
}

func (detector *stringDetector) Name() string {
	return "string"
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	ruleReferenceType settings.RuleReferenceType,
	evaluator types.Evaluator,
) ([]interface{}, error) {
	if node.Type() == "string_literal" {
		return []interface{}{generictypes.String{
			Value:     node.Content(),
			IsLiteral: true,
		}}, nil
	}

	return nil, nil
}
func (detector *stringDetector) Close() {}
