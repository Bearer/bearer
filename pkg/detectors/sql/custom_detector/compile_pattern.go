package customdetector

import (
	"regexp"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/detectors/custom/config"
	"github.com/bearer/bearer/pkg/parser"
	"github.com/bearer/bearer/pkg/parser/custom"
	"github.com/bearer/bearer/pkg/parser/nodeid"
	"github.com/bearer/bearer/pkg/util/file"
)

var tableNameRegex = regexp.MustCompile(`\$TABLE_NAME`)
var columnsRegex = regexp.MustCompile(`[<]?\$COLUMN[>]?`)
var functionNameRegex = regexp.MustCompile(`\$FUNCTION_NAME`)
var scriptRegex = regexp.MustCompile(`\$SCRIPT`)
var anythingRegex = regexp.MustCompile(`\$ANTYHING`)

func (detector *Detector) CompilePattern(
	rulePattern settings.RulePattern,
	idGenerator nodeid.Generator,
) (config.CompiledRule, error) {
	reworkedRule := tableNameRegex.ReplaceAllLiteral(
		[]byte(rulePattern.Pattern),
		[]byte("Var_Table_Name"+idGenerator.GenerateId()),
	)
	reworkedRule = columnsRegex.ReplaceAllLiteral(reworkedRule, []byte("Var_Column_Id"+idGenerator.GenerateId()+" Var_Type_id"+idGenerator.GenerateId()))
	reworkedRule = functionNameRegex.ReplaceAllLiteral(reworkedRule, []byte("Var_Function_Name"+idGenerator.GenerateId()+"()"))
	reworkedRule = scriptRegex.ReplaceAllLiteral(reworkedRule, []byte(`$$ Var_Script`+idGenerator.GenerateId()+` $$`))
	reworkedRule = anythingRegex.ReplaceAllLiteral(reworkedRule, []byte(`Var_Anything`+idGenerator.GenerateId()))

	tree, err := parser.ParseBytes(&file.FileInfo{}, &file.Path{}, []byte(reworkedRule), language, 0)
	if err != nil {
		return config.CompiledRule{}, err
	}
	defer tree.Close()

	compiledRule := &config.CompiledRule{
		Params: make([]config.Param, 0),
	}

	custom.GenerateTreeSitterQuery(tree.RootNode().Child(0), idGenerator, compiledRule, detector, false)

	return *compiledRule, nil
}

func (detector *Detector) Annotate(tree *parser.Tree) error {
	return nil
}
