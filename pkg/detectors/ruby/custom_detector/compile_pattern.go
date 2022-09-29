package customdetector

import (
	"regexp"

	"github.com/bearer/curio/pkg/detectors/custom/config"
	"github.com/bearer/curio/pkg/parser"
	"github.com/bearer/curio/pkg/parser/custom"
	"github.com/bearer/curio/pkg/parser/nodeid"
	"github.com/bearer/curio/pkg/util/file"
)

var classNameRegex = regexp.MustCompile(`\$CLASS_NAME`)
var argumentsRegex = regexp.MustCompile(`<\$ARGUMENT>`)
var anythingRegex = regexp.MustCompile(`\$ANYTHING`)

func (detector *Detector) CompilePattern(Rule string, idGenerator nodeid.Generator) (config.CompiledRule, error) {
	reworkedRule := classNameRegex.ReplaceAll([]byte(Rule), []byte("Var_Class_Name"+idGenerator.GenerateId()))
	reworkedRule = argumentsRegex.ReplaceAll([]byte(reworkedRule), []byte("Var_Arguments"+idGenerator.GenerateId()))
	reworkedRule = anythingRegex.ReplaceAll([]byte(reworkedRule), []byte("Var_Anything"+idGenerator.GenerateId()))

	tree, err := parser.ParseBytes(&file.FileInfo{}, &file.Path{}, []byte(reworkedRule), language, 0)
	if err != nil {
		return config.CompiledRule{}, err
	}
	defer tree.Close()

	rule := &config.CompiledRule{
		Languages: []string{"ruby"},
		Params:    make([]config.Param, 0),
	}

	custom.GenerateTreeSitterQuery(tree.RootNode().Child(0), idGenerator, rule, detector)

	return *rule, nil
}
