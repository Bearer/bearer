package string

import (
	"fmt"
	"regexp"

	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/stringutil"

	"github.com/bearer/bearer/internal/scanner/detectors/common"
	"github.com/bearer/bearer/internal/scanner/detectors/types"
)

var (
	simpleInterpolationRegexp  = regexp.MustCompile(`%#?[a-zA-Z]`)
	numericInterpolationRegexp = regexp.MustCompile(`%\d?\.?\d?f`)
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "call_expression":
		function := node.ChildByFieldName("function")
		if function.Type() == "selector_expression" {
			field := function.ChildByFieldName("field")
			if field.Type() == "field_identifier" && field.Content() == "Sprintf" {
				arguments := node.ChildByFieldName("arguments").NamedChildren()

				stringValue, isLiteral, err := common.GetStringValue(arguments[0], detectorContext)
				if err != nil || !isLiteral {
					return nil, err
				}

				stringValue = simpleInterpolationRegexp.ReplaceAllString(stringValue, "%s")  // %s %d %#v %t
				stringValue = numericInterpolationRegexp.ReplaceAllString(stringValue, "%s") // %2.2f %.2f %2f %2.f

				newArguments := []any{}
				for index, argument := range arguments {
					if index == 0 {
						continue
					}

					childValue, childIsLiteral, err := common.GetStringValue(argument, detectorContext)
					if err != nil {
						return nil, err
					}

					if !childIsLiteral {
						isLiteral = false
						if childValue == "" {
							childValue = common.NonLiteralValue
						}
					}

					newArguments = append(newArguments, childValue)
				}

				value := fmt.Sprintf(stringValue, newArguments...)

				return []interface{}{common.String{
					Value:     value,
					IsLiteral: isLiteral,
				}}, nil
			}
		}
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "assignment_statement":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	case "interpreted_string_literal", "raw_string_literal":
		value := stringutil.StripQuotes(node.Content())

		return []interface{}{common.String{
			Value:     value,
			IsLiteral: true,
		}}, nil
	}

	return nil, nil
}
