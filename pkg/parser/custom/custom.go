package custom

import (
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/detectors/custom/config"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/report/schema/datatype"
)

type Detector interface {
	ExtractArguments(node *parser.Node, idGenerator nodeid.Generator) (map[parser.NodeID]*datatype.DataType, error)
	CompilePattern(rulePattern settings.RulePattern, idGenerator nodeid.Generator) (config.CompiledRule, error)
	IsParam(node *parser.Node) (bool, bool, *config.Param)
	Annotate(tree *parser.Tree) error
}

func GenerateTreeSitterQuery(node *parser.Node, idGenerator nodeid.Generator, rule *config.CompiledRule, detector Detector, isChild bool) {
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

	if param != nil && param.MatchAnything {
		rule.Tree += "_"
	} else {
		rule.Tree += node.Type() + ""
	}

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

		GenerateTreeSitterQuery(child, idGenerator, rule, detector, true)
	}

	rule.Tree += ")"

	if !isChild {
		rule.Tree += " @rule"
	}
}
