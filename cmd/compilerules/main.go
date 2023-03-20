package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bearer/bearer/pkg/ast/languages/ruby"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/maputil"
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

	var rules []*settings.Rule
	for _, rule := range maputil.ToSortedSlice(config.BuiltInRules) {
		if slices.Contains(rule.Languages, "ruby") {
			rules = append(rules, rule)
		}
	}
	for _, rule := range maputil.ToSortedSlice(config.Rules) {
		if slices.Contains(rule.Languages, "ruby") {
			rules = append(rules, rule)
		}
	}

	language.WriteRules(rules, writer)

	ruleFile.Close()
	return compiler.Compile("souffle/rules.dl", "pkg/souffle/rules/generated.cpp")
}
