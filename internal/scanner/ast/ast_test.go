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

type ruleInfo struct {
	ID    string
	Index int
}

func TestDisabledRules(t *testing.T) {
	content := `
		# bearer:disable rule1,rule2
		# bearer:disable rule3
		def m(a)
			# bearer:disable rule4
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
			"rule4": {Id: "rule4", Languages: languageIDs},
		},
	)
	if err != nil {
		t.Fatalf("failed to create rule set: %s", err)
	}

	var ruleDump []ruleInfo
	for _, rule := range ruleSet.Rules() {
		if rule.Type() != ruleset.RuleTypeBuiltin {
			ruleDump = append(ruleDump, ruleInfo{ID: rule.ID(), Index: rule.Index()})
		}
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
		ruleDump,
		tree.RootNode().Dump(),
	)
}

func TestExpectedRules(t *testing.T) {
	content := `
		# bearer:expected rule1
		def m(a)
			b.bar
		end
	`

	language := ruby.Get()
	languageIDs := []string{language.ID()}

	ruleSet, err := ruleset.New(
		language.ID(),
		map[string]*settings.Rule{
			"rule1": {Id: "rule1", Languages: languageIDs},
		},
	)
	if err != nil {
		t.Fatalf("failed to create rule set: %s", err)
	}

	var ruleDump []ruleInfo
	for _, rule := range ruleSet.Rules() {
		if rule.Type() != ruleset.RuleTypeBuiltin {
			ruleDump = append(ruleDump, ruleInfo{ID: rule.ID(), Index: rule.Index()})
		}
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
		ruleDump,
		tree.RootNode().Dump(),
	)
}
