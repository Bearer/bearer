package language

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
)

type Language interface {
	Name() string
	EnryLanguages() []string
	NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector
	SitterLanguage() *sitter.Language
	Pattern() Pattern
	NewAnalyzer(builder *tree.Builder) Analyzer
}

type Analyzer interface {
	Analyze(node *sitter.Node, visitChildren func() error) error
}
