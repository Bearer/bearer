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
	NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector
	SitterLanguage() *sitter.Language
	Pattern() Pattern
	NewAnalyzer(builder *tree.Builder) Analyzer
}

type Analyzer interface {
	Analyze(node *sitter.Node, visitChildren func() error) error
}
