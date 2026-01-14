package language

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
)

type Language interface {
	ID() string
	DisplayName() string
	EnryLanguages() []string
	GoclocLanguages() []string
	NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector
	SitterLanguage() *sitter.Language
	Pattern() Pattern
	NewAnalyzer(builder *tree.Builder) Analyzer
	// StringFragmentTypes returns the node types that represent literal string content
	// for this language. These types are used by EachContentPart to identify text content
	// vs interpolated/dynamic content within string nodes.
	// Examples: "string_fragment" (JS/Java), "string_content" (Python/Ruby), "string_value" (PHP)
	StringFragmentTypes() []string
}

type Analyzer interface {
	Analyze(node *sitter.Node, visitChildren func() error) error
}
