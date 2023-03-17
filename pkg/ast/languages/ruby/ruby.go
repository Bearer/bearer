package ruby

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/bearer/pkg/ast/languages/ruby/patterns"
	"github.com/bearer/bearer/pkg/ast/sourcefacts"
	"github.com/bearer/bearer/pkg/ast/walker"
	"github.com/bearer/bearer/pkg/souffle/writer"
	filewriter "github.com/bearer/bearer/pkg/souffle/writer/file"
)

type Language struct {
	sitterLanguage *sitter.Language
	walker         *walker.Walker
}

func New() *Language {
	sitterLanguage := ruby.GetLanguage()

	return &Language{
		sitterLanguage: sitterLanguage,
		walker:         walker.NewWalker(sitterLanguage),
	}
}

func (language *Language) WriteSourceFacts(
	fileId uint32,
	input []byte,
	writer writer.FactWriter,
) error {
	rootNode, err := sitter.ParseCtx(context.TODO(), input, language.sitterLanguage)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := sourcefacts.WriteFacts(
		language.walker,
		fileId,
		input,
		rootNode,
		writer,
	); err != nil {
		return fmt.Errorf("translation error: %w", err)
	}

	return nil
}

func (language *Language) WriteRule(
	ruleName string,
	input string,
	writer *filewriter.Writer,
) error {
	inputBytes := []byte(input)

	rootNode, err := sitter.ParseCtx(context.TODO(), inputBytes, language.sitterLanguage)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	if err := patterns.CompileRule(
		language.walker,
		ruleName,
		inputBytes,
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
