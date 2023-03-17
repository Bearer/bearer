package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bearer/bearer/pkg/ast/languages/ruby"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/souffle/compiler"
	filewriter "github.com/bearer/bearer/pkg/util/souffle/writer/file"
	"golang.org/x/exp/slices"
)

func main() {
	if err := do(); err != nil {
		log.Fatal(err)
	}
}

func do() error {
	language := ruby.New()
	defer language.Close()

	ruleFile, err := os.Create("souffle/generated/rules.dl")
	if err != nil {
		return fmt.Errorf("error creating souffle rules file: %w", err)
	}
	defer ruleFile.Close()

	config, err := settings.FromOptions(flag.Options{})
	if err != nil {
		return fmt.Errorf("error loading rules: %w", err)
	}

	writer := filewriter.New(ruleFile)

	for ruleName, rule := range config.BuiltInRules {
		if err := compileRuleToSouffle(language, writer, ruleName, rule); err != nil {
			return fmt.Errorf("built-in rule %s error: %w", ruleName, err)
		}
	}

	for ruleName, rule := range config.Rules {
		if err := compileRuleToSouffle(language, writer, ruleName, rule); err != nil {
			return fmt.Errorf("rule %s error: %w", ruleName, err)
		}
	}

	ruleFile.Close()
	return compiler.Compile("souffle/rules.dl", "pkg/souffle/rules/generated.cpp")
}

func compileRuleToSouffle(
	language *ruby.Language,
	writer *filewriter.Writer,
	ruleName string,
	rule *settings.Rule,
) error {
	if !slices.Contains(rule.Languages, "ruby") {
		return nil
	}

	for _, pattern := range rule.Patterns {
		if err := language.WriteRule(ruleName, pattern.Pattern, writer); err != nil {
			return fmt.Errorf("pattern error (%s)': %w", pattern.Pattern, err)
		}
	}

	return nil
}
