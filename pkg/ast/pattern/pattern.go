package pattern

import (
	"fmt"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/patternquery/builder"
	languagetypes "github.com/bearer/bearer/new/language/types"
	treetypes "github.com/bearer/bearer/pkg/ast/tree/types"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	sitter "github.com/smacker/go-tree-sitter"
)

type Variable struct {
	name string
}

type Pattern struct {
	query          *sitter.Query
	variables      []Variable
	nameToVariable map[string]*Variable
}

type Match struct {
	node      *sitter.Node
	variables []*sitter.Node
}

func Parse(
	lang languagetypes.Language,
	langImplementation implementation.Implementation,
	content string,
	focusedVariable string,
	filters []settings.PatternFilter,
) (*Pattern, error) {
	builderResult, err := builder.Build(lang, langImplementation, content, focusedVariable)
	if err != nil {
		return nil, fmt.Errorf("error building sitter query: %w", err)
	}

	query, err := sitter.NewQuery([]byte(builderResult.Query), langImplementation.SitterLanguage())
	if err != nil {
		return nil, fmt.Errorf("error compiling sitter query: %w", err)
	}

	variables := make([]Variable, len(builderResult.VariableNames))
	for i, name := range builderResult.VariableNames {
		variables[i].name = name
	}

	nameToVariable := make(map[string]*Variable)
	for i, variable := range variables {
		nameToVariable[variable.name] = &variables[i]
	}

	return &Pattern{
		query:          query,
		variables:      variables,
		nameToVariable: nameToVariable,
	}, nil
}

func (pattern *Pattern) InsertMatches(treeBuilder *treetypes.Builder, sitterRootNode *sitter.Node) {

}

func (pattern *Pattern) Close() {
	pattern.query.Close()
}
