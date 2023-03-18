package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bearer/bearer/pkg/ast/languages/ruby"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/souffle/compiler"
	writerbase "github.com/bearer/bearer/pkg/util/souffle/writer/base"
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
		if !slices.Contains(rule.Languages, "ruby") {
			continue
		}

		if err := compileRuleToSouffle(language, writer, ruleName, rule); err != nil {
			return fmt.Errorf("built-in rule %s error: %w", ruleName, err)
		}
	}

	for ruleName, rule := range config.Rules {
		if !slices.Contains(rule.Languages, "ruby") {
			continue
		}

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
	ruleRelation := fmt.Sprintf("Rule_Pattern_%s", ruleName)
	writer.WriteRelation(ruleRelation, "patternIndex: Rule_PatternIndex", "node: AST_NodeId")

	variableRelation := fmt.Sprintf("Rule_PatternVariable_%s", ruleName)
	// writer.WriteRelation(variableRelation, "node: AST_NodeId", "variable: Rule_VariableName", "variableNode: AST_NodeId")

	if err := writer.WriteRule(
		[]writerbase.Predicate{writer.Predicate(
			"Rule",
			writer.Symbol(ruleName),
			writer.Identifier("patternIndex"),
			writer.Identifier("node"),
		)},
		[]writerbase.Literal{writer.Predicate(
			ruleRelation,
			writer.Identifier("patternIndex"),
			writer.Identifier("node"),
		)},
	); err != nil {
		return fmt.Errorf("error writing generic rule: %w", err)
	}

	for i, pattern := range rule.Patterns {
		if err := language.WriteRule(ruleRelation, variableRelation, i, pattern.Pattern, writer); err != nil {
			return fmt.Errorf("pattern error (%s)': %w", pattern.Pattern, err)
		}
	}

	return nil
}
