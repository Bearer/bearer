package javascript

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"

	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/javascript/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/javascript/string"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/javascript/analyzer"
	"github.com/bearer/bearer/new/language/implementation/javascript/pattern"
)

type javascriptImplementation struct {
	pattern pattern.Pattern
}

func Get() implementation.Implementation {
	return &javascriptImplementation{}
}

func (*javascriptImplementation) Name() string {
	return "javascript"
}

func (*javascriptImplementation) EnryLanguages() []string {
	return []string{"JavaScript", "TypeScript", "TSX"}
}

func (*javascriptImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJavascript, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*javascriptImplementation) SitterLanguage() *sitter.Language {
	return tsx.GetLanguage()
}

func (implementation *javascriptImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*javascriptImplementation) NewAnalyzer(builder *tree.Builder) implementation.Analyzer {
	return analyzer.New(builder)
}
