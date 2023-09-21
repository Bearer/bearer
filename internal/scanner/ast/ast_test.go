package ast_test

import (
	"context"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/languages/ruby"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

func TestDisabledRules(t *testing.T) {
	content := `
		# bearer:disable rule1
		# bearer:disable rule2
		def m(a)
			# bearer:disable rule3
			a.foo
			b.bar
		end
	`

	language := ruby.Get()
	languageIDs := []string{language.ID()}

	ruleSet, err := ruleset.New(
		language.ID(),
		map[string]*settings.Rule{
			"rule1": {Id: "rule1", Languages: languageIDs},
			"rule2": {Id: "rule2", Languages: languageIDs},
			"rule3": {Id: "rule3", Languages: languageIDs},
		},
	)
	if err != nil {
		t.Fatalf("failed to create rule set: %s", err)
	}

	querySet := query.NewSet(language.ID(), language.SitterLanguage())
	if err := querySet.Compile(); err != nil {
		t.Fatalf("failed to compile query set: %s", err)
	}

	tree, err := ast.ParseAndAnalyze(
		context.Background(),
		language,
		ruleSet,
		querySet,
		[]byte(content),
	)

	if err != nil {
		t.Fatalf("failed to parse and analyze input: %s", err)
	}

	cupaloy.SnapshotT(
		t,
		ruleSet.Rules(),
		tree.RootNode().Dump(),
	)
}
