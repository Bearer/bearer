package ruby

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	sitterruby "github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/ruby"
	builderinput "github.com/bearer/bearer/new/language/patternquery/builder/input"
	"github.com/bearer/bearer/pkg/ast/idgenerator"
	"github.com/bearer/bearer/pkg/ast/languages/ruby/patterns"
	"github.com/bearer/bearer/pkg/ast/sourcefacts"
	"github.com/bearer/bearer/pkg/ast/walker"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/souffle/writer"
	filewriter "github.com/bearer/bearer/pkg/util/souffle/writer/file"
)

type Language struct {
	sitterLanguage     *sitter.Language
	langImplementation implementation.Implementation
	walker             *walker.Walker
}

func New() *Language {
	sitterLanguage := sitterruby.GetLanguage()

	return &Language{
		sitterLanguage:     sitterLanguage,
		langImplementation: ruby.Get(),
		walker:             walker.NewWalker(sitterLanguage),
	}
}

func (language *Language) WriteSourceFacts(
	input []byte,
	rootNode *sitter.Node,
	nodeIdGenerator *idgenerator.NodeIdGenerator,
	writer writer.FactWriter,
) error {
	return sourcefacts.WriteFacts(
		language.walker,
		input,
		rootNode,
		nodeIdGenerator,
		writer,
	)
}

func (language *Language) WriteRules(rules []*settings.Rule, writer *filewriter.Writer) error {
	for _, rule := range rules {
		writer.WriteComment(fmt.Sprintf("rule %s", rule.Id))

		for patternIndex, pattern := range rule.Patterns {
			patternId := idgenerator.PatternId(rule.Id, patternIndex)

			if err := language.WriteRule(patternId, pattern.Pattern, writer); err != nil {
				return fmt.Errorf("pattern error (%s)': %w", pattern.Pattern, err)
			}
		}
	}

	return nil
}

func (language *Language) WriteRule(patternId string, input string, writer *filewriter.Writer) error {
	processedInput, inputParams, err := builderinput.Process(language.langImplementation, input)
	if err != nil {
		return fmt.Errorf("error parsing bearer syntax: %w", err)
	}

	processedInputBytes := []byte(processedInput)

	rootNode, err := sitter.ParseCtx(context.TODO(), processedInputBytes, language.sitterLanguage)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := patterns.CompileRule(
		language.walker,
		inputParams,
		patternId,
		processedInputBytes,
		rootNode,
		writer,
	); err != nil {
		return fmt.Errorf("error compiling rule: %w", err)
	}

	return nil
}

func (language *Language) Close() {
	language.walker.Close()
}
