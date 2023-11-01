package php_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/languages/php"
	"github.com/bearer/bearer/internal/languages/testhelper"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	patternquerybuilder "github.com/bearer/bearer/internal/scanner/detectors/customrule/patternquery/builder"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

//go:embed testdata/logger.yml
var loggerRule []byte

//go:embed testdata/scope_rule.yml
var scopeRule []byte

func TestFlow(t *testing.T) {
	testhelper.GetRunner(t, loggerRule, "PHP").RunTest(t, "./testdata/testcases/flow", ".snapshots/flow/")
}

func TestScope(t *testing.T) {
	testhelper.GetRunner(t, scopeRule, "PHP").RunTest(t, "./testdata/scope", ".snapshots/")
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
