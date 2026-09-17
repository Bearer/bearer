package php_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/go-enry/go-enry/v2"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/languages/php"
	"github.com/bearer/bearer/pkg/languages/testhelper"
	"github.com/bearer/bearer/pkg/scanner/ast"
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	patternquerybuilder "github.com/bearer/bearer/pkg/scanner/detectors/customrule/patternquery/builder"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

//go:embed testdata/md.yml
var mdRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, php.Get()).RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, php.Get()).RunTest(t, "./testdata/scope", ".snapshots/")
}

func TestConst(t *testing.T) {
	testhelper.GetRunner(t, mdRule, php.Get()).RunTest(t, "./testdata/md", ".snapshots/")
}

func TestAnalyzer(t *testing.T) {
	for _, test := range []struct{ name, code string }{
		{"foreach", `<?php
		    $array = [];

				foreach ($array as $value) {
					echo $value;
				}

				foreach ($array as $key => $value) {
					echo $key;
					echo $value;
				}
		`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			language := php.Get()

			ruleSet, err := ruleset.New(language.ID(), make(map[string]*settings.Rule))
			if err != nil {
				tt.Fatalf("failed to create rule set: %s", err)
			}

			querySet := query.NewSet(language.ID(), language.SitterLanguage())
			if err := querySet.Compile(); err != nil {
				tt.Fatalf("failed to compile query set: %s", err)
			}

			result, err := ast.ParseAndAnalyze(context.Background(), language, ruleSet, querySet, []byte(test.code))
			if err != nil {
				tt.Fatalf("failed to parse example: %s", err)
			}

			cupaloy.SnapshotT(tt, result.RootNode().Dump())
		})
	}
}

func TestPattern(t *testing.T) {
	for _, test := range []struct{ name, pattern string }{
		{"class name in object creation is unanchored", `
				new $<!>Foo;
		`},
		{"named arguments are unanchored", `
				foo(x: $<!>$<_>)
		`},
		{"property names are unanchored", `
				class $<_> {
					public $<!>$<_>;
				}
		`},
		{"parameter names are unanchored", `
				class $<_> {
					public function $<_>($<_> $<!>$<_>) {}
				}
		`},
		{"catch clauses and types are unanchored", `
				try {} catch ($<_> $<!>$$<_>) {}
		`},
		{"function names and bodies are unanchored", `
				function $<_>() {}
		`},
		{"anonymous function names and bodies are unanchored", `
				function () {};
		`},
	} {
		t.Run(test.name, func(tt *testing.T) {
			result, err := patternquerybuilder.Build(php.Get(), test.pattern, "")
			if err != nil {
				tt.Fatalf("failed to build pattern: %s", err)
			}

			cupaloy.SnapshotT(tt, result)
		})
	}
}

// enry classifies .phtml as HTML+PHP rather than PHP, and the scanner skips a
// file whose detected language is not in this list. .phtml is the default
// template extension for Magento 2 and Laminas, so leaving it out means those
// templates are never scanned.
func TestEnryLanguagesCoversPhtml(t *testing.T) {
	languages := php.Get().EnryLanguages()

	assert.Contains(t, languages, "PHP")
	assert.Contains(t, languages, "HTML+PHP")

	detected, _ := enry.GetLanguageByExtension("template.phtml")
	assert.Equal(t, "HTML+PHP", detected, "enry classification changed; this list needs to follow it")
	assert.Contains(t, languages, detected)
}
