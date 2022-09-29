package custom

import (
	"github.com/bearer/curio/pkg/detectors/custom/config"
	"github.com/bearer/curio/pkg/parser"
	parserdatatype "github.com/bearer/curio/pkg/parser/datatype"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

type Detector interface {
	ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*parserdatatype.DataType, error)
	CompilePattern(Rule string, idGenerator nodeid.Generator) (config.CompiledRule, error)
	IsParam(node *parser.Node) (bool, bool, *config.Param)
}

func GenerateTreeSitterQuery(node *parser.Node, idGenerator nodeid.Generator, rule *config.CompiledRule, detector Detector) {
	if node.Type() == "ERROR" {
		return
	}

	end, shouldIgnore, param := detector.IsParam(node)

	if shouldIgnore {
		return
	}

	varName := "var"
	paramVar := "param_" + varName

	if rule.Tree != "" {
		rule.Tree += " "
	}
	rule.Tree += "("
	rule.Tree += node.Type() + ""

	assignedID := ""

	if param != nil {
		assignedID = idGenerator.GenerateId()
		paramName := paramVar + assignedID
		rule.Tree += ")"
		rule.Tree += " @" + paramName
		param.ParamName = varName + assignedID
		rule.Params = append(rule.Params, *param)
		if end {
			return
		}
	}

	for i := 0; i < node.ChildCount(); i++ {
		child := node.Child(i)

		// this approach fails for complicated queries such as Rails.log.info($ARGUMENTS)
		// startString := node.Sitter().String()
		// // remove our params as they don't exist in tree
		// trimStart := string(regexpTreeSitterParams.ReplaceAll([]byte(rule.Tree), []byte("")))
		// endMatch := node.Child(i).Sitter().String()

		// // remove start string
		// cutset := strings.TrimPrefix(startString, trimStart)
		// endPosition := strings.Index(cutset, endMatch)

		// if endPosition > 0 {
		// 	fieldName := strings.Trim(cutset[0:endPosition], " ")
		// 	if fieldName != "" {
		// 		rule.Tree += " " + fieldName
		// 	}
		// }

		// latest library version 7621c203ae43fe58c0fc4d18e4dcf3caa1985888 supports this but it is buggy so its of no use
		// childFieldName := node.FieldNameForChild(i)
		// if childFieldName != "" {
		// 	rule.Tree += " " + childFieldName + ": "
		// }

		GenerateTreeSitterQuery(child, idGenerator, rule, detector)
	}

	rule.Tree += ")"
}
